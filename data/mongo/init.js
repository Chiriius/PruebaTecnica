db = db.getSiblingDB('events_db');

db.events.insertMany([
    {
        name: "Conferencia de Tecnología",
        type: "Conferencia",
        description: "Evento anual de tecnología e innovación",
        date: new Date(),
        status: "Pendiente por revisar"
    },
    {
        name: "Taller de Go",
        type: "Taller",
        description: "Taller práctico de programación en Go",
        date: new Date(Date.now() + 86400000), 
        status: "Pendiente por revisar"
    },
    {
        name: "Reunión de Equipo",
        type: "Reunión",
        description: "Reunión mensual del equipo de desarrollo",
        date: new Date(Date.now() - 86400000), 
        status: "Revisado",
        category: "Sin gestión"
    },
    {
        name: "Incidente de Seguridad",
        type: "Incidente",
        description: "Reporte de posible vulnerabilidad",
        date: new Date(),
        status: "Revisado",
        category: "Requiere gestión"
    }
]);
