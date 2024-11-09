package econewscron

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	openai "github.com/sashabaranov/go-openai"
)

func filterUnparsedNews(news []*EconomicNew, db *sqlx.DB) []*EconomicNew {
	urls := make([]string, len(news))
	for i := range news {
		urls[i] = news[i].URL
	}
	sql := `
	SELECT ARRAY(
		SELECT EXISTS (
			SELECT 1 FROM economic_new WHERE url = u.url
		) FROM unnest($1::text[]) as u(url)
	);`

	exists := make([]bool, len(news))
	err := db.Get(pq.Array(&exists), sql, pq.Array(urls))
	if err != nil {
		log.Printf("Error quering existing news: %v", err)
		return make([]*EconomicNew, 0)
	}

	filteredNews := make([]*EconomicNew, 0)
	for i := range news {
		if !exists[i] {
			filteredNews = append(filteredNews, news[i])
		}
	}

	return filteredNews
}

func saveEcoNews(db *sqlx.DB, news []*EconomicNew) {
	if len(news) == 0 {
		return
	}
	dbNews := make([]*DbEconomicNew, len(news))
	for i := range news {
		dbNews[i] = news[i].ToDbEconomicNew()
		dbNews[i].CreatedAt = time.Now().UTC()
	}
	sql := `
		INSERT INTO economic_new (id, title, date, url, image, summary, company, tags, sentiment, created_at, deleted_at)
		VALUES (:id, :title, :date, :url, :image, :summary, :company, :tags, :sentiment, :created_at, :deleted_at);
	`
	_, err := db.NamedExec(sql, dbNews)
	if err != nil {
		log.Printf("Error saving economic news: %v", err)
		return
	}
}

const openAIEconomicNewPromptOld = `
Eres un economista especializado en la economía boliviana. Se te proporcionará una noticia económica, el título esta entre las etiquetas <title></title> y el contenido entre las etiquetas <content></content>. Debes realizar los siguientes pasos:

1. Determinar si la noticia tiene un impacto directo a la economía boliviana, es decir, si tiene relación  MUY directa con el precio de los productos de la canasta familiar, productos mayor consumidos, precio del dólar o un incremento o decremento de las fuentes de empleo. Ten en cuenta sólo noticias que traten de hecho o sucesos ya acontecidos o confirmados. Descartar noticias que traten de acontecimientos futuros, rumores, anuncios u opiniones. En caso de que la noticia sea descartada Resumen debe ser null, Sentimiento 0 y Tags un array vacío.

2. Si la noticia no se descartó, generar un resumen, sentimiento y tags de acuerto  a las siguientes indicaciones.

Resumen: Generar un encabezado claro y conciso de la noticia de 40 a 60 palabras, dado que es un encabezado, puedes omitir detalles si no se puede resumir  en menos de 60 palabras, puedes usar etiquetas HTML de texto para resaltar o dar mejor formato (es decir puedes usar <i> <strong> <b> etc.)

Sentimiento: Un número entre 1 y  10 denotanto el impacto negativo o positivo de la noticia, 1 es extremadamente negativo y 10 es extremadamente positivo, ten en cuenta que 5 se considera noticia negativa y 6 positiva.

Tags: Un array con las palabras clave de la noticia para una posterior clasificación y  búsqueda de relacionados, la longitud puede ser desde 1 a 6 tags. Ignora tags como Bolivia, económia, porque todas las noticias son sobre la economía boliviana.
`

const openAIEconomicNewPrompt = `Eres un economista especializado en la economía boliviana, con un criterio muy estricto para evaluar el impacto directo de noticias económicas. Se te proporcionará una noticia económica, el título esta entre las etiquetas <title></title> y el contenido entre las etiquetas <content></content>. Debes realizar los siguientes pasos:

Determinar si la noticia tiene un impacto DIRECTO E INMEDIATO en la economía boliviana. La noticia SOLO debe considerarse relevante si cumple AL MENOS UNO de estos criterios específicos:
a) Afecta de manera inmediata y cuantificable los precios de productos básicos de la canasta familiar (ejemplo: incremento confirmado en el precio del pan, arroz, aceite, etc.)
b) Modifica el precio del dólar en más de 1% en un solo día
c) Impacta directamente el empleo formal con cifras concretas (ejemplo: cierre confirmado de empresas grandes, despidos masivos documentados, apertura de nuevas fuentes de trabajo con números específicos)
d) Cambios confirmados en salarios o beneficios que afecten a más del 10% de la población trabajadora
e) Medidas económicas gubernamentales YA IMPLEMENTADAS (no anuncios) que modifiquen precios, impuestos o aranceles de productos de consumo masivo

IMPORTANTE: Descartar AUTOMÁTICAMENTE:

Noticias sobre planes futuros o intenciones
Anuncios o promesas sin implementación
Opiniones o análisis de expertos
Rumores o especulaciones de mercado
Noticias sobre sectores que no afectan directamente al consumidor promedio
Estadísticas o informes que no implican cambios inmediatos
Noticias sobre negociaciones en curso
Proyecciones económicas

En caso de que la noticia sea descartada, Resumen debe ser null, Sentimiento 0 y Tags un array vacío.

Si la noticia cumple estrictamente con al menos uno de los criterios anteriores, generar:

Resumen: Generar un encabezado claro y conciso de la noticia de 40 a 60 palabras, enfocándose específicamente en el impacto económico directo. Usar etiquetas HTML de texto para resaltar cifras y datos concretos (es decir puedes usar <i> <strong> <b> etc.)
Sentimiento: Un número entre 1 y 10 denotando el impacto negativo o positivo de la noticia:
1-2: Impacto catastrófico en la economía familiar
3-4: Impacto muy negativo en la economía familiar
5: Impacto moderadamente negativo
6: Impacto moderadamente positivo
7-8: Impacto muy positivo en la economía familiar
9-10: Impacto extraordinariamente positivo
Tags: Un array con las palabras clave específicas de la noticia (1 a 4 tags), enfocándose en el tipo de impacto económico (ejemplo: 'precios', 'empleo', 'salarios', 'dólar', etc.). Ignorar tags genéricos como 'Bolivia', 'economía'.
`

type jsonSchema map[string]interface{}

func (s jsonSchema) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}(s))
}

func generateAIEcoNewData(ecoNew *EconomicNew, openAIKey string) {
	type response struct {
		Resumen     *string  `json:"resumen"`
		Sentimiento int      `json:"sentimiento"`
		Tags        []string `json:"tags"`
	}
	userMessage := fmt.Sprintf("<title>%s</title><content>%s</content>", ecoNew.Title, *ecoNew.Content)
	jsonSchema := jsonSchema{
		"type": "object",
		"properties": map[string]interface{}{
			"resumen":     map[string]string{"type": "string"},
			"sentimiento": map[string]string{"type": "number"},
			"tags": map[string]interface{}{
				"type":  "array",
				"items": map[string]string{"type": "string"},
			},
		},
		"additionalProperties": false,
		"required": []string{
			"resumen",
			"sentimiento",
			"tags",
		},
	}

	client := openai.NewClient(openAIKey)
	AIresp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4o20240806,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: openAIEconomicNewPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userMessage,
				},
			},
			ResponseFormat: &openai.ChatCompletionResponseFormat{
				Type: openai.ChatCompletionResponseFormatTypeJSONSchema,
				JSONSchema: &openai.ChatCompletionResponseFormatJSONSchema{
					Name:   "eco_new_response",
					Strict: true,
					Schema: jsonSchema,
				},
			},
		},
	)
	if err != nil {
		log.Printf("Error generating AI eco new data for %s: %v", ecoNew.ID, err)
	}
	respJSONMessage := AIresp.Choices[0].Message.Content
	var resp response
	err = json.Unmarshal([]byte(respJSONMessage), &resp)
	if err != nil {
		log.Printf("Error unmarshaling AI response for %s: %v", ecoNew.ID, err)
	}
	ecoNew.Summary = resp.Resumen
	ecoNew.Sentiment = int8(resp.Sentimiento)
	ecoNew.Tags = resp.Tags
}
