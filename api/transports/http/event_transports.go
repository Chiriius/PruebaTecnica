package transports

import (
	"context"
	"net/http"
	"prueba_tecnica/api/endpoints"
	"prueba_tecnica/api/entities"
	"prueba_tecnica/api/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @title			API Prueba tecnica Gestión de Eventos
// @version		1.0
// @description	API para la gestión y clasificación de eventos
// @host			localhost:8080
// @BasePath		/api/v1
func NewEventRouter(router *gin.Engine, endpoints endpoints.EventEndpoints, logger logrus.FieldLogger) {
	eventGroup := router.Group("/api/v1/events")

	//	@Summary		Crear un nuevo evento
	//	@Description	Crea un nuevo evento en el sistema
	//	@Tags			Eventos
	//	@Accept			json
	//	@Produce		json
	//	@Param			event	body		entities.Event		true	"Datos del Evento"
	//	@Success		201		{object}	map[string]string	"ID del evento creado"
	//	@Failure		400		{object}	map[string]string	"Error en los datos de entrada"
	//	@Failure		500		{object}	map[string]string	"Error interno del servidor"
	//	@Router			/events [post]
	eventGroup.POST("/", func(c *gin.Context) {
		var event entities.Event
		if err := c.ShouldBindJSON(&event); err != nil {
			logger.Errorln("Layer:event_transports", "Method: Post", "Error:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		transportEvent, err := endpoints.CreateEvent(context.Background(), event)

		if err == service.ErrStatus {
			logger.Errorln("Layer:event_transports", "Method: Post", "Error:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err == service.ErrValidation {
			logger.Errorln("Layer:event_transports", "Method: Post", "Error:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err != nil {
			logger.Errorln("Layer:event_transports", "Method: Post", "Error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		logger.Infoln("Layer:event_transports", "Method: Post", "Event:", transportEvent.ID)
		c.JSON(http.StatusCreated, gin.H{"id": transportEvent.ID})
	})

	//	@Summary		Obtener un evento por ID
	//	@Description	Obtiene los detalles de un evento específico
	//	@Tags			Eventos
	//	@Produce		json
	//	@Param			id	path		string				true	"ID del Evento"
	//	@Success		200	{object}	entities.Event		"Evento encontrado"
	//	@Failure		404	{object}	map[string]string	"Evento no encontrado"
	//	@Failure		500	{object}	map[string]string	"Error interno del servidor"
	//	@Router			/events/{id} [get]
	eventGroup.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		event, err := endpoints.GetEventByID(c.Request.Context(), id)
		if err != nil {
			logger.Errorln("Layer:event_transports", "Method: GET", "Error:", err)
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		logger.Infoln("Layer:event_transports", "Method: GET", "Event:", event)
		c.JSON(http.StatusOK, event)
	})

	//	@Summary		Listar todos los eventos
	//	@Description	Obtiene una lista de todos los eventos registrados
	//	@Tags			Eventos
	//	@Produce		json
	//	@Success		200	{array}		entities.Event		"Lista de eventos"
	//	@Failure		500	{object}	map[string]string	"Error interno del servidor"
	//	@Router			/events [get]
	eventGroup.GET("/", func(c *gin.Context) {
		events, err := endpoints.GetAllEvents(c.Request.Context())
		if err != nil {
			logger.Errorln("Layer:event_transports", "Method: GET", "Error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener eventos: " + err.Error()})
			return
		}
		logger.Infoln("Layer:event_transports", "Method: GET", "Eventos: obtenidos correctamente")
		c.JSON(http.StatusOK, events)
	})

	//	@Summary		Actualizar un evento
	//	@Description	Actualiza los datos de un evento existente
	//	@Tags			Eventos
	//	@Accept			json
	//	@Produce		json
	//	@Param			id		path		string				true	"ID del Evento"
	//	@Param			event	body		entities.Event		true	"Datos actualizados del Evento"
	//	@Success		200		{object}	entities.Event		"Evento actualizado"
	//	@Failure		400		{object}	map[string]string	"Error en los datos de entrada"
	//	@Failure		404		{object}	map[string]string	"Evento no encontrado"
	//	@Failure		500		{object}	map[string]string	"Error interno del servidor"
	//	@Router			/events/{id} [put]
	eventGroup.PUT("/:id", func(c *gin.Context) {
		id := c.Param("id")
		var event entities.Event
		if err := c.ShouldBindJSON(&event); err != nil {
			logger.Errorln("Layer:event_transports", "Method: PUT", "Error:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos del evento inválidos: " + err.Error()})
			return
		}
		event.ID = id
		_, err := endpoints.UpdateEvent(c.Request.Context(), event)

		if err == service.ErrValidation {
			logger.Errorln("Layer:event_transports", "Method: PUT", "Error:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err == service.ErrStatus {
			logger.Errorln("Layer:event_transports", "Method: PUT", "Error:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err == service.ErrEventNotfound {
			logger.Errorln("Layer:event_transports", "Method: PUT", "Error:", err)
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if err != nil {
			logger.Errorln("Layer:event_transports", "Method: PUT", "Error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		logger.Infoln("Layer:event_transports", "Method: PUT", "Eventos: Actualizado correctamente")
		c.JSON(http.StatusOK, event)
	})

	//	@Summary		Eliminar un evento
	//	@Description	Elimina un evento existente del sistema
	//	@Tags			Eventos
	//	@Produce		json
	//	@Param			id	path		string				true	"ID del Evento"
	//	@Success		200	{object}	map[string]string	"Mensaje de confirmación"
	//	@Failure		400	{object}	map[string]string	"Error en la solicitud"
	//	@Failure		500	{object}	map[string]string	"Error interno del servidor"
	//	@Router			/events/{id} [delete]
	eventGroup.DELETE("/:id", func(c *gin.Context) {
		id := c.Param("id")
		if err := endpoints.DeleteEvent(c.Request.Context(), id); err != nil {
			logger.Errorln("Layer:event_transports", "Method: DELETE", "Error:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		logger.Infoln("Layer:event_transports", "Method: DELETE", "Evento: Eliminado correctamente")
		c.JSON(http.StatusOK, gin.H{"message": "Evento eliminado correctamente"})
	})

	//	@Summary		Clasificar evento automáticamente
	//	@Description	Clasifica automáticamente un evento revisado según su tipo
	//	@Tags			Clasificación
	//	@Produce		json
	//	@Param			id	path		string				true	"ID del Evento"
	//	@Success		200	{object}	map[string]string	"Mensaje de confirmación"
	//	@Failure		400	{object}	map[string]string	"Error en la solicitud"
	//	@Failure		500	{object}	map[string]string	"Error interno del servidor"
	//	@Router			/events/{id}/classify [put]
	eventGroup.PUT("/:id/classify", func(c *gin.Context) {
		id := c.Param("id")
		if _, err := endpoints.ClassifyEvent(c.Request.Context(), id); err != nil {
			logger.Errorln("Layer:event_transports", "Method: PUT", "Error:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error al clasificar evento: " + err.Error()})
			return
		}
		logger.Infoln("Layer:event_transports", "Method: PUT", "Evento:clasificado automáticamente según su tipo")
		c.JSON(http.StatusOK, gin.H{"message": "Evento clasificado automáticamente según su tipo"})
	})

	//	@Summary		Clasificar evento manualmente
	//	@Description	Permite clasificar manualmente un evento revisado
	//	@Tags			Clasificación
	//	@Accept			json
	//	@Produce		json
	//	@Param			id			path		string				true	"ID del Evento"
	//	@Param			category	body		string				true	"Categoría ('Requiere gestión' o 'Sin gestión')"
	//	@Success		200			{object}	map[string]string	"Mensaje de confirmación"
	//	@Failure		400			{object}	map[string]string	"Error en la solicitud"
	//	@Failure		500			{object}	map[string]string	"Error interno del servidor"
	//	@Router			/events/{id}/manual-classify [put]
	eventGroup.PUT("/:id/manual-classify", func(c *gin.Context) {
		id := c.Param("id")
		var request struct {
			Category string `json:"category" binding:"required,oneof='Requiere gestión' 'Sin gestión'"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			logger.Errorln("Layer:event_transports", "Method: PUT", "Error:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Categoría inválida: " + err.Error()})
			return
		}
		if _, err := endpoints.ManualClassifyEvent(c.Request.Context(), id, request.Category); err != nil {
			logger.Errorln("Layer:event_transports", "Method: PUT", "Error:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error en clasificación manual: " + err.Error()})
			return
		}
		logger.Infoln("Layer:event_transports", "Method: PUT", "Evento:clasificado manualmente")
		c.JSON(http.StatusOK, gin.H{"message": "Evento clasificado manualmente"})
	})

	//	@Summary		Filtrar eventos por estado
	//	@Description	Obtiene una lista de eventos filtrados por estado
	//	@Tags			Consultas
	//	@Produce		json
	//	@Param			status	path		string				true	"Estado del evento (Pendiente/Revisado)"
	//	@Success		200		{array}		entities.Event		"Lista de eventos filtrados"
	//	@Failure		400		{object}	map[string]string	"Error en la solicitud"
	//	@Failure		500		{object}	map[string]string	"Error interno del servidor"
	//	@Router			/events/status/{status} [get]
	eventGroup.GET("/status/:status", func(c *gin.Context) {
		status := c.Param("status")
		events, err := endpoints.GetEventsByStatus(c.Request.Context(), status)
		if err != nil {
			logger.Errorln("Layer:event_transports", "Method: GET", "Error:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error al filtrar por estado: " + err.Error()})
			return
		}
		logger.Infoln("Layer:event_transports", "Method: GET", "Eventos por estado obtenidos correctamente")
		c.JSON(http.StatusOK, events)
	})

	//	@Summary		Filtrar eventos por categoría
	//	@Description	Obtiene una lista de eventos filtrados por categoría
	//	@Tags			Consultas
	//	@Produce		json
	//	@Param			category	path		string				true	"Categoría del evento"
	//	@Success		200			{array}		entities.Event		"Lista de eventos filtrados"
	//	@Failure		400			{object}	map[string]string	"Error en la solicitud"
	//	@Failure		500			{object}	map[string]string	"Error interno del servidor"
	//	@Router			/events/category/{category} [get]
	eventGroup.GET("/category/:category", func(c *gin.Context) {
		category := c.Param("category")
		events, err := endpoints.GetEventsByCategory(c.Request.Context(), category)
		if err != nil {
			logger.Errorln("Layer:event_transports", "Method: GET", "Error:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error al filtrar por categoría: " + err.Error()})
			return
		}
		logger.Infoln("Layer:event_transports", "Method: GET", "Eventos por categoria obtenidos correctamente")
		c.JSON(http.StatusOK, events)
	})

	// Obtener eventos que requieren gestión
	//	@Summary		Obtener eventos que requieren gestión
	//	@Description	Obtiene una lista de eventos marcados como que requieren gestión
	//	@Tags			Consultas
	//	@Produce		json
	//	@Success		200	{array}		entities.Event		"Lista de eventos que requieren gestión"
	//	@Failure		500	{object}	map[string]string	"Error interno del servidor"
	//	@Router			/events/needs [get]
	eventGroup.GET("/needs", func(c *gin.Context) {
		events, err := endpoints.GetEventsNeedingAction(c.Request.Context())
		if err != nil {
			logger.Errorln("Layer:event_transports", "Method: GET", "Error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener eventos: " + err.Error()})
			return
		}
		logger.Infoln("Layer:event_transports", "Method: GET", "Eventos que requieren gestion obtenidos correctamente")
		c.JSON(http.StatusOK, events)
	})
}
