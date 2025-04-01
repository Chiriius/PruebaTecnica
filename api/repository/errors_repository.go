package repository

import "errors"

var ErrEventNotfound = errors.New("Error evento no encontrado")
var ErrNotasks = errors.New("No sé elimino ningun evento, ese evento no se encuentre en la base de datos")
