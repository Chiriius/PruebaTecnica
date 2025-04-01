package service

import "errors"

var ErrValidation = errors.New("Error en la estructura del request llene todos los campos")
var ErrStatus = errors.New("El estado debe de ser Pendiente por revisar o Revisado)")
var ErrEventNotfound = errors.New("evento con ese id no encontrado")
var ErrTypeCategory = errors.New("categoría inválida")
var ErrNoID = errors.New("Id del evento requerido")
var ErrCategory = errors.New("categoría debe ser 'Requiere gestión' o 'Sin gestión'")
var ErrEventRevi = errors.New("Solo se pueden clasificar eventos revisados")
