package seed

import (
	"log"
	"time"

	"github.com/isaacunaa/ticketek-ds2026/backend/internal/domain"
	"github.com/isaacunaa/ticketek-ds2026/backend/internal/utils"
	"gorm.io/gorm"
)

func Ejecutar(db *gorm.DB) {
	seedUsuarios(db)
	seedEventos(db)
}

func seedUsuarios(db *gorm.DB) {
	var count int64
	db.Model(&domain.Usuario{}).Count(&count)
	if count > 0 {
		return
	}

	usuarios := []struct {
		email    string
		password string
		nombre   string
		apellido string
		rol      string
	}{
		{"admin@vortice.com", "admin123", "Admin", "Vórtice", "admin"},
		{"juan@mail.com", "123456", "Juan", "Pérez", "cliente"},
		{"maria@mail.com", "123456", "María", "González", "cliente"},
	}

	for _, u := range usuarios {
		hash, err := utils.HashearPassword(u.password)
		if err != nil {
			log.Printf("Error hasheando password de %s: %v", u.email, err)
			continue
		}
		usuario := domain.Usuario{
			Email:        u.email,
			PasswordHash: hash,
			Nombre:       u.nombre,
			Apellido:     u.apellido,
			Rol:          u.rol,
		}
		if err := db.Create(&usuario).Error; err != nil {
			log.Printf("Error creando usuario %s: %v", u.email, err)
		}
	}
	log.Println("Seed: usuarios creados")
}

func seedEventos(db *gorm.DB) {
	eventos := []domain.Evento{
		// ── Música ──────────────────────────────────────────────
		{
			Titulo:          "Lollapalooza Argentina 2026",
			Descripcion:     "El festival de música más grande de Latinoamérica vuelve a Buenos Aires con los mejores artistas del mundo.",
			FechaHora:       time.Date(2026, 3, 20, 14, 0, 0, 0, time.UTC),
			DuracionMinutos: 480,
			Ubicacion:       "Hipódromo de San Isidro, Buenos Aires",
			Categoria:       "música",
			CupoTotal:       50000,
			CupoDisponible:  50000,
			Precio:          25000,
			Estado:          "activo",
			ImagenURL:       "https://picsum.photos/seed/lollapalooza/600/400",
		},
		{
			Titulo:          "Cosquín Rock 2026",
			Descripcion:     "El festival de rock nacional e internacional más importante del interior del país.",
			FechaHora:       time.Date(2026, 2, 7, 16, 0, 0, 0, time.UTC),
			DuracionMinutos: 600,
			Ubicacion:       "Santa María de Punilla, Córdoba",
			Categoria:       "música",
			CupoTotal:       30000,
			CupoDisponible:  30000,
			Precio:          18000,
			Estado:          "activo",
			ImagenURL:       "https://picsum.photos/seed/cosquinrock/600/400",
		},
		{
			Titulo:          "Noche de Tango en el Colón",
			Descripcion:     "Una velada única de tango en el Teatro Colón con las mejores orquestas del país.",
			FechaHora:       time.Date(2026, 7, 12, 21, 0, 0, 0, time.UTC),
			DuracionMinutos: 120,
			Ubicacion:       "Teatro Colón, Buenos Aires",
			Categoria:       "música",
			CupoTotal:       2500,
			CupoDisponible:  2500,
			Precio:          12000,
			Estado:          "activo",
			ImagenURL:       "https://picsum.photos/seed/tango/600/400",
		},
		{
			Titulo:          "Festival de Jazz Buenos Aires",
			Descripcion:     "Cuatro días de jazz en vivo con artistas nacionales e internacionales al aire libre.",
			FechaHora:       time.Date(2026, 8, 20, 18, 0, 0, 0, time.UTC),
			DuracionMinutos: 240,
			Ubicacion:       "Usina del Arte, Buenos Aires",
			Categoria:       "música",
			CupoTotal:       5000,
			CupoDisponible:  5000,
			Precio:          0,
			Estado:          "activo",
			ImagenURL:       "https://picsum.photos/seed/jazz/600/400",
		},
		// ── Teatro ──────────────────────────────────────────────
		{
			Titulo:          "El Método Gronholm",
			Descripcion:     "Cuatro candidatos compiten por un puesto ejecutivo en una dinámica de selección que se convierte en una trampa.",
			FechaHora:       time.Date(2026, 5, 8, 20, 30, 0, 0, time.UTC),
			DuracionMinutos: 100,
			Ubicacion:       "Teatro San Martín, Buenos Aires",
			Categoria:       "teatro",
			CupoTotal:       800,
			CupoDisponible:  800,
			Precio:          7500,
			Estado:          "activo",
			ImagenURL:       "https://picsum.photos/seed/gronholm/600/400",
		},
		{
			Titulo:          "Chicago: El Musical",
			Descripcion:     "El clásico de Broadway llega a Buenos Aires con producción internacional.",
			FechaHora:       time.Date(2026, 9, 5, 21, 0, 0, 0, time.UTC),
			DuracionMinutos: 150,
			Ubicacion:       "Teatro Ópera, Buenos Aires",
			Categoria:       "teatro",
			CupoTotal:       1500,
			CupoDisponible:  1500,
			Precio:          14000,
			Estado:          "activo",
			ImagenURL:       "https://picsum.photos/seed/chicago/600/400",
		},
		{
			Titulo:          "Hamlet — Compañía Nacional",
			Descripcion:     "La tragedia más universal de Shakespeare en una puesta en escena contemporánea.",
			FechaHora:       time.Date(2026, 6, 25, 20, 0, 0, 0, time.UTC),
			DuracionMinutos: 130,
			Ubicacion:       "Teatro Cervantes, Buenos Aires",
			Categoria:       "teatro",
			CupoTotal:       600,
			CupoDisponible:  600,
			Precio:          6000,
			Estado:          "activo",
			ImagenURL:       "https://picsum.photos/seed/hamlet/600/400",
		},
		// ── Humor ───────────────────────────────────────────────
		{
			Titulo:          "Stand Up: Noche de Comedia",
			Descripcion:     "Los mejores comediantes del país en una noche inolvidable.",
			FechaHora:       time.Date(2026, 5, 15, 21, 0, 0, 0, time.UTC),
			DuracionMinutos: 90,
			Ubicacion:       "Teatro Gran Rex, Buenos Aires",
			Categoria:       "humor",
			CupoTotal:       1200,
			CupoDisponible:  1200,
			Precio:          8000,
			Estado:          "activo",
			ImagenURL:       "https://picsum.photos/seed/standup/600/400",
		},
		{
			Titulo:          "Fede Bal: Tour 2026",
			Descripcion:     "El actor y comediante recorre el país con su nuevo show de stand up.",
			FechaHora:       time.Date(2026, 4, 18, 21, 30, 0, 0, time.UTC),
			DuracionMinutos: 90,
			Ubicacion:       "Luna Park, Buenos Aires",
			Categoria:       "humor",
			CupoTotal:       8000,
			CupoDisponible:  8000,
			Precio:          9500,
			Estado:          "activo",
			ImagenURL:       "https://picsum.photos/seed/fedebal/600/400",
		},
		{
			Titulo:          "Improv Fest — Festival de Improvisación",
			Descripcion:     "Doce compañías de improvisación teatral se enfrentan en un torneo de risas garantizadas.",
			FechaHora:       time.Date(2026, 10, 3, 20, 0, 0, 0, time.UTC),
			DuracionMinutos: 120,
			Ubicacion:       "El Galpón de Guevara, Buenos Aires",
			Categoria:       "humor",
			CupoTotal:       350,
			CupoDisponible:  350,
			Precio:          4500,
			Estado:          "activo",
			ImagenURL:       "https://picsum.photos/seed/improv/600/400",
		},
		// ── Deportes ────────────────────────────────────────────
		{
			Titulo:          "River vs Boca — Superclásico",
			Descripcion:     "El partido más esperado del año en el Monumental.",
			FechaHora:       time.Date(2026, 4, 5, 20, 0, 0, 0, time.UTC),
			DuracionMinutos: 120,
			Ubicacion:       "Estadio Monumental, Buenos Aires",
			Categoria:       "deportes",
			CupoTotal:       8000,
			CupoDisponible:  8000,
			Precio:          15000,
			Estado:          "activo",
			ImagenURL:       "https://picsum.photos/seed/superclasico/600/400",
		},
		{
			Titulo:          "Maratón de Buenos Aires 2026",
			Descripcion:     "La carrera más popular de Argentina con categorías de 21k y 42k.",
			FechaHora:       time.Date(2026, 10, 18, 7, 0, 0, 0, time.UTC),
			DuracionMinutos: 300,
			Ubicacion:       "Palermo, Buenos Aires",
			Categoria:       "deportes",
			CupoTotal:       15000,
			CupoDisponible:  15000,
			Precio:          3500,
			Estado:          "activo",
			ImagenURL:       "https://picsum.photos/seed/maraton/600/400",
		},
		{
			Titulo:          "Argentina vs Brasil — Amistoso FIFA",
			Descripcion:     "La Selección recibe a Brasil en el Monumental en partido amistoso de preparación.",
			FechaHora:       time.Date(2026, 11, 14, 21, 0, 0, 0, time.UTC),
			DuracionMinutos: 120,
			Ubicacion:       "Estadio Monumental, Buenos Aires",
			Categoria:       "deportes",
			CupoTotal:       70000,
			CupoDisponible:  70000,
			Precio:          22000,
			Estado:          "activo",
			ImagenURL:       "https://picsum.photos/seed/argbrasil/600/400",
		},
		// ── Arte ────────────────────────────────────────────────
		{
			Titulo:          "arteBA 2026 — Feria de Arte Contemporáneo",
			Descripcion:     "La feria de arte más importante de América Latina con más de 100 galerías.",
			FechaHora:       time.Date(2026, 5, 21, 14, 0, 0, 0, time.UTC),
			DuracionMinutos: 360,
			Ubicacion:       "La Rural, Palermo, Buenos Aires",
			Categoria:       "arte",
			CupoTotal:       10000,
			CupoDisponible:  10000,
			Precio:          2500,
			Estado:          "activo",
			ImagenURL:       "https://picsum.photos/seed/arteba/600/400",
		},
		{
			Titulo:          "Van Gogh Alive — Experiencia Inmersiva",
			Descripcion:     "La exposición multisensorial más visitada del mundo llega a Buenos Aires.",
			FechaHora:       time.Date(2026, 7, 1, 10, 0, 0, 0, time.UTC),
			DuracionMinutos: 60,
			Ubicacion:       "Centro Cultural Kirchner, Buenos Aires",
			Categoria:       "arte",
			CupoTotal:       500,
			CupoDisponible:  500,
			Precio:          5500,
			Estado:          "activo",
			ImagenURL:       "https://picsum.photos/seed/vangogh/600/400",
		},
		{
			Titulo:          "Noche de los Museos 2026",
			Descripcion:     "Una noche mágica con todos los museos de la ciudad abiertos y con actividades especiales.",
			FechaHora:       time.Date(2026, 11, 7, 19, 0, 0, 0, time.UTC),
			DuracionMinutos: 480,
			Ubicacion:       "Ciudad de Buenos Aires",
			Categoria:       "arte",
			CupoTotal:       999999,
			CupoDisponible:  999999,
			Precio:          0,
			Estado:          "activo",
			ImagenURL:       "https://picsum.photos/seed/museos/600/400",
		},
		// ── Espectáculo ─────────────────────────────────────────
		{
			Titulo:          "Cirque du Soleil — Alegría",
			Descripcion:     "El espectáculo de circo contemporáneo más famoso del mundo.",
			FechaHora:       time.Date(2026, 6, 10, 20, 30, 0, 0, time.UTC),
			DuracionMinutos: 120,
			Ubicacion:       "Parque de la Ciudad, Buenos Aires",
			Categoria:       "espectáculo",
			CupoTotal:       3000,
			CupoDisponible:  3000,
			Precio:          18000,
			Estado:          "activo",
			ImagenURL:       "https://picsum.photos/seed/cirque/600/400",
		},
		{
			Titulo:          "Fuerza Bruta — Wayra",
			Descripcion:     "La experiencia teatral inmersiva argentina que conquista el mundo vuelve a casa.",
			FechaHora:       time.Date(2026, 8, 14, 21, 0, 0, 0, time.UTC),
			DuracionMinutos: 70,
			Ubicacion:       "Galpón Fuerza Bruta, Palermo",
			Categoria:       "espectáculo",
			CupoTotal:       600,
			CupoDisponible:  600,
			Precio:          11000,
			Estado:          "activo",
			ImagenURL:       "https://picsum.photos/seed/fuerzabruta/600/400",
		},
		{
			Titulo:          "Disney on Ice — Frozen",
			Descripcion:     "Los personajes más queridos de Disney cobran vida sobre el hielo.",
			FechaHora:       time.Date(2026, 7, 25, 17, 0, 0, 0, time.UTC),
			DuracionMinutos: 90,
			Ubicacion:       "Arena Ciudad de Buenos Aires",
			Categoria:       "espectáculo",
			CupoTotal:       12000,
			CupoDisponible:  12000,
			Precio:          8500,
			Estado:          "activo",
			ImagenURL:       "https://picsum.photos/seed/disney/600/400",
		},
		// ── Cine ────────────────────────────────────────────────
		{
			Titulo:          "Cine Bajo las Estrellas — Clásicos del 90",
			Descripcion:     "Ciclo de proyecciones al aire libre de las películas más queridas de los años 90.",
			FechaHora:       time.Date(2026, 1, 30, 21, 30, 0, 0, time.UTC),
			DuracionMinutos: 120,
			Ubicacion:       "Bosques de Palermo, Buenos Aires",
			Categoria:       "cine",
			CupoTotal:       1000,
			CupoDisponible:  1000,
			Precio:          2000,
			Estado:          "activo",
			ImagenURL:       "https://picsum.photos/seed/cineestrellas/600/400",
		},
		{
			Titulo:          "Festival Internacional de Cine de Mar del Plata",
			Descripcion:     "El festival de cine más importante de América Latina con más de 200 películas.",
			FechaHora:       time.Date(2026, 11, 6, 10, 0, 0, 0, time.UTC),
			DuracionMinutos: 600,
			Ubicacion:       "Mar del Plata, Buenos Aires",
			Categoria:       "cine",
			CupoTotal:       5000,
			CupoDisponible:  5000,
			Precio:          3000,
			Estado:          "activo",
			ImagenURL:       "https://picsum.photos/seed/cinemardelplata/600/400",
		},
		// ── Tecnología ──────────────────────────────────────────
		{
			Titulo:          "Nerdearla 2026",
			Descripcion:     "La conferencia de tecnología más grande de Argentina con charlas, talleres y networking.",
			FechaHora:       time.Date(2026, 10, 8, 9, 0, 0, 0, time.UTC),
			DuracionMinutos: 480,
			Ubicacion:       "Teatro Ópera, Buenos Aires",
			Categoria:       "tecnología",
			CupoTotal:       3000,
			CupoDisponible:  3000,
			Precio:          5000,
			Estado:          "activo",
			ImagenURL:       "https://picsum.photos/seed/nerdearla/600/400",
		},
		{
			Titulo:          "AI Summit Argentina 2026",
			Descripcion:     "La cumbre de inteligencia artificial más importante del país con líderes de la industria.",
			FechaHora:       time.Date(2026, 9, 18, 9, 0, 0, 0, time.UTC),
			DuracionMinutos: 480,
			Ubicacion:       "Centro Costa Salguero, Buenos Aires",
			Categoria:       "tecnología",
			CupoTotal:       2000,
			CupoDisponible:  2000,
			Precio:          12000,
			Estado:          "activo",
			ImagenURL:       "https://picsum.photos/seed/aisummit/600/400",
		},
	}

	agregados := 0
	actualizados := 0
	for _, e := range eventos {
		var existing domain.Evento
		result := db.Where("titulo = ?", e.Titulo).First(&existing)
		if result.Error != nil {
			// No existe, crear
			if err := db.Create(&e).Error; err != nil {
				log.Printf("Error creando evento '%s': %v", e.Titulo, err)
			} else {
				agregados++
			}
		} else if e.ImagenURL != "" {
			// Siempre sincronizar imagen_url
			if err := db.Model(&existing).Update("imagen_url", e.ImagenURL).Error; err != nil {
				log.Printf("Error actualizando imagen de '%s': %v", e.Titulo, err)
			} else {
				actualizados++
			}
		}
	}
	if agregados > 0 {
		log.Printf("Seed: %d evento(s) nuevos creados", agregados)
	}
	if actualizados > 0 {
		log.Printf("Seed: %d evento(s) con imagen actualizada", actualizados)
	}

	// Fix directo para eventos con títulos que contienen em-dash (coincidencia LIKE)
	fixes := []struct {
		like string
		url  string
	}{
		{"Cirque%", "https://picsum.photos/seed/cirque/600/400"},
		{"River%Boca%", "https://picsum.photos/seed/superclasico/600/400"},
	}
	for _, f := range fixes {
		db.Model(&domain.Evento{}).Where("titulo LIKE ?", f.like).Update("imagen_url", f.url)
	}
	log.Println("Seed: fix de imágenes aplicado")
}
