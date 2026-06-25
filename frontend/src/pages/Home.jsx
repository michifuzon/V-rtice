import { useState, useEffect, useCallback, useRef } from 'react'
import { Link } from 'react-router-dom'
import client from '../api/client'

const token = () => localStorage.getItem('token')

const EMOJI = {
  'música':      '🎵',
  'humor':       '😂',
  'teatro':      '🎭',
  'deportes':    '⚽',
  'arte':        '🎨',
  'cine':        '🎬',
  'tecnología':  '💻',
  'espectáculo': '🎪',
}
const getEmoji = (cat) => {
  if (!cat) return '🎉'
  const cats = cat.split(',').map(c => c.trim().toLowerCase())
  for (const c of cats) { if (EMOJI[c]) return EMOJI[c] }
  return '🎉'
}

const CATEGORIAS = [
  { label: 'Todos',        value: 'Todos' },
  { label: 'Música',       value: 'música' },
  { label: 'Humor',        value: 'humor' },
  { label: 'Teatro',       value: 'teatro' },
  { label: 'Deportes',     value: 'deportes' },
  { label: 'Arte',         value: 'arte' },
  { label: 'Espectáculo',  value: 'espectáculo' },
]

/* ── Modal cancelar membresía ────────────────────────────── */
function ModalCancelarClub({ onClose, onConfirm, cancelando }) {
  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-box" onClick={e => e.stopPropagation()}>
        <div className="modal-header">
          <h3>Cancelar membresía</h3>
          <button className="modal-close-btn" onClick={onClose}>✕</button>
        </div>
        <div className="modal-body">
          <p className="modal-desc">
            ¿Confirmás la cancelación de tu membresía del <strong>Club Vórtice</strong>?
            Perderás el descuento del 10% en todas las entradas de forma inmediata.
          </p>
          <div className="modal-footer">
            <button className="btn-modal-cancel" onClick={onClose}>No, mantener</button>
            <button className="btn-modal-confirm btn-modal-danger" onClick={onConfirm} disabled={cancelando}>
              {cancelando ? 'Cancelando...' : 'Sí, cancelar membresía'}
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}

/* ── Club Vórtice Banner ─────────────────────────────────── */
function ClubVortice({ suscripcion, onSuscribir, suscribiendo, onCancelar, cancelando }) {
  if (!token()) return null

  const activa = suscripcion?.activa
  const vence  = suscripcion?.vence
    ? new Date(suscripcion.vence).toLocaleDateString('es-AR', { day: 'numeric', month: 'long', year: 'numeric' })
    : null

  if (activa) {
    return (
      <div className="club-banner club-banner-activo">
        <div className="club-banner-inner">
          <div className="club-banner-info">
            <div className="club-activo-header">
              <span className="club-activo-check">✓</span>
              <span className="club-activo-titulo">Sos miembro del Club Vórtice</span>
            </div>
            <p className="club-activo-detalle">
              Tu descuento del <strong>10%</strong> está activo en todas las entradas
              {vence && <> · Membresía válida hasta el <strong>{vence}</strong></>}
            </p>
          </div>
          <div className="club-activo-acciones">
            <button className="club-btn club-btn-renovar" onClick={onSuscribir} disabled={suscribiendo || cancelando}>
              {suscribiendo ? 'Procesando...' : '↺ Renovar'}
            </button>
            <button className="club-btn-cancelar" onClick={onCancelar} disabled={cancelando || suscribiendo}>
              {cancelando ? 'Cancelando...' : 'Cancelar membresía'}
            </button>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="club-banner">
      <div className="club-banner-inner">
        <div className="club-banner-info">
          <div className="club-promo-header">
            <span className="club-badge">💜 Club Vórtice</span>
            <span className="club-promo-precio">$5.000 / mes</span>
          </div>
          <p className="club-desc">
            Unite al club y conseguí un <strong>10% de descuento</strong> en <strong>todas las entradas</strong> del catálogo.
          </p>
          <div className="club-beneficios">
            <span>✓ Descuento inmediato</span>
            <span>✓ Sin renovación automática</span>
            <span>✓ Cancelable cuando quieras</span>
          </div>
        </div>
        <button className="club-btn" onClick={onSuscribir} disabled={suscribiendo}>
          {suscribiendo ? 'Procesando...' : 'Suscribirme — $5.000/mes'}
        </button>
      </div>
    </div>
  )
}

/* ── EventCard ───────────────────────────────────────────── */
function EventCard({ evento, esFavorito, onToggleFavorito, suscripcionActiva }) {
  const nombre = evento.nombre || evento.titulo || 'Sin nombre'
  const emoji  = getEmoji(evento.categoria)
  const [imgError, setImgError] = useState(false)

  const fecha = evento.fecha_hora
    ? new Date(evento.fecha_hora).toLocaleDateString('es-AR', { day: 'numeric', month: 'short', year: 'numeric' })
    : null

  const hayDescuento = suscripcionActiva && evento.precio > 0
  const precioDesc   = hayDescuento ? Math.round(evento.precio * 0.9 * 100) / 100 : null

  const precioLabel =
    evento.precio == null ? null :
    evento.precio === 0   ? 'Gratis' :
    `$${Number(evento.precio).toLocaleString('es-AR')}`

  return (
    <div className="event-card-wrap">
      <Link to={`/eventos/${evento.id}`} className="event-card">
        {evento.imagen_url && !imgError
          ? <img src={evento.imagen_url} alt={nombre} className="event-img" onError={() => setImgError(true)} />
          : <div className="event-placeholder">{emoji}</div>
        }

        <div className="event-body">
          {evento.categoria && (
            <div className="event-badges">
              {evento.categoria.split(',').map(c => c.trim()).map(c => (
                <span key={c} className="event-category-badge">{c}</span>
              ))}
            </div>
          )}
          <h3 className="event-title">{nombre}</h3>
          {evento.ubicacion && <p className="event-venue">📍 {evento.ubicacion}</p>}
          {fecha             && <p className="event-date">📅 {fecha}</p>}
          {precioLabel && (
            <div style={{ marginTop: '.3rem' }}>
              {hayDescuento ? (
                <div className="precio-descuento-wrap">
                  <span className="precio-tachado">{precioLabel}</span>
                  <span className="precio-nuevo">${Number(precioDesc).toLocaleString('es-AR')}</span>
                  <span className="badge-descuento">10% OFF</span>
                </div>
              ) : (
                <p style={{ fontWeight: 700, color: 'var(--accent-dark)', fontSize: '.92rem' }}>{precioLabel}</p>
              )}
            </div>
          )}
        </div>

        <div className="event-footer">
          <span className="btn-buy">Ver más</span>
        </div>
      </Link>

      {token() && (
        <button
          className={`heart-btn${esFavorito ? ' heart-active' : ''}`}
          onClick={() => onToggleFavorito?.(evento.id)}
          title={esFavorito ? 'Quitar de favoritos' : 'Agregar a favoritos'}
        >
          {esFavorito ? '♥' : '♡'}
        </button>
      )}
    </div>
  )
}

/* ── Home ────────────────────────────────────────────────── */
export default function Home() {
  const [eventos,   setEventos]   = useState([])
  const [search,    setSearch]    = useState('')
  const [categoria, setCategoria] = useState(CATEGORIAS[0])
  const [loading,   setLoading]   = useState(true)
  const [error,     setError]     = useState('')

  const [favoritoIds,   setFavoritoIds]   = useState(new Set())
  const [suscripcion,   setSuscripcion]   = useState(null)
  const [suscribiendo,  setSuscribiendo]  = useState(false)
  const [cancelando,       setCancelando]       = useState(false)
  const [modalCancelarClub, setModalCancelarClub] = useState(false)
  const [clubFeedback,  setClubFeedback]  = useState('')

  const inputRef = useRef(null)

  useEffect(() => {
    if (!token()) return
    client.get('/favoritos')
      .then(({ data }) => setFavoritoIds(new Set((data.favoritos ?? []).map(fav => fav.evento_id))))
      .catch(() => {})
    client.get('/club/estado')
      .then(({ data }) => setSuscripcion(data))
      .catch(() => {})
  }, [])

  const toggleFavorito = async (eventoId) => {
    if (!token()) return
    try {
      if (favoritoIds.has(eventoId)) {
        await client.delete(`/favoritos/${eventoId}`)
        setFavoritoIds(prev => { const n = new Set(prev); n.delete(eventoId); return n })
      } else {
        await client.post(`/favoritos/${eventoId}`)
        setFavoritoIds(prev => new Set([...prev, eventoId]))
      }
    } catch { /* ignore */ }
  }

  const handleSuscribir = async () => {
    setSuscribiendo(true)
    setClubFeedback('')
    try {
      const { data } = await client.post('/club/suscribir')
      // Re-consultar el estado real desde el servidor para asegurar consistencia
      const { data: estado } = await client.get('/club/estado')
      setSuscripcion(estado)
      setClubFeedback(data.mensaje || '¡Suscripción activada! Ahora tenés 10% de descuento en todas las entradas.')
      setTimeout(() => setClubFeedback(''), 5000)
    } catch (err) {
      const status = err.response?.status
      const msg =
        err.response?.data?.error ||
        (status === 404 ? 'Reiniciá el backend para activar el Club Vórtice.' :
         status === 401 ? 'Sesión expirada, iniciá sesión nuevamente.' :
         status >= 500  ? 'Error en el servidor. Verificá que el backend esté actualizado.' :
         !err.response  ? 'No se pudo conectar con el servidor.' :
                          'No se pudo procesar la suscripción.')
      setClubFeedback('❌ ' + msg)
    } finally {
      setSuscribiendo(false)
    }
  }

  const fetchEventos = useCallback(async (searchVal, catVal) => {
    setLoading(true)
    setError('')
    try {
      const params = {}
      if (searchVal.trim()) params.search = searchVal.trim()

      const { data } = await client.get('/eventos', { params })
      let lista = Array.isArray(data) ? data : data.eventos ?? data.data ?? []

      if (catVal !== 'Todos') {
        lista = lista.filter(ev => {
          if (!ev.categoria) return false
          return ev.categoria.split(',').map(c => c.trim()).includes(catVal)
        })
      }

      setEventos(lista)
    } catch (err) {
      setError('No se pudieron cargar los eventos. Verificá tu conexión e intentá de nuevo.')
      console.error(err)
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    fetchEventos(search, categoria.value)
  }, [categoria.value]) // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    const t = setTimeout(() => { fetchEventos(search, categoria.value) }, 400)
    return () => clearTimeout(t)
  }, [search]) // eslint-disable-line react-hooks/exhaustive-deps

  const handleCancelarSuscripcion = async () => {
    setCancelando(true)
    try {
      const { data } = await client.delete('/club/suscripcion')
      const { data: estado } = await client.get('/club/estado')
      setSuscripcion(estado)
      setModalCancelarClub(false)
      setClubFeedback(data.mensaje || 'Membresía cancelada.')
      setTimeout(() => setClubFeedback(''), 5000)
    } catch {
      setClubFeedback('❌ No se pudo cancelar la membresía.')
      setModalCancelarClub(false)
    } finally {
      setCancelando(false)
    }
  }

  const handleClearFilters = () => {
    setSearch('')
    setCategoria(CATEGORIAS[0])
    if (inputRef.current) inputRef.current.value = ''
  }

  return (
    <>
      {/* ── Hero ─────────────────────────────────────────── */}
      <section className="hero">
        <h1 className="hero-title">Encontrá tu próximo evento</h1>
        <p className="hero-subtitle">Eventos que atrapan</p>

        <div className="hero-search">
          <input
            ref={inputRef}
            type="text"
            className="hero-search-input"
            placeholder="Buscar eventos..."
            defaultValue={search}
            onChange={(e) => setSearch(e.target.value)}
          />
          <span className="hero-search-icon">🔍</span>
        </div>
      </section>

      {/* ── Club Vórtice ─────────────────────────────────── */}
      <ClubVortice
        suscripcion={suscripcion}
        onSuscribir={handleSuscribir}
        suscribiendo={suscribiendo}
        onCancelar={() => setModalCancelarClub(true)}
        cancelando={cancelando}
      />

      {clubFeedback && (
        <div className={`club-feedback${clubFeedback.startsWith('❌') ? ' club-feedback-error' : ''}`}>
          {clubFeedback}
          <button onClick={() => setClubFeedback('')} style={{ background: 'none', border: 'none', cursor: 'pointer', marginLeft: '.5rem', opacity: .7 }}>✕</button>
        </div>
      )}

      {/* ── Barra de categorías ───────────────────────────── */}
      <div className="category-bar">
        {CATEGORIAS.map((cat) => (
          <button
            key={cat.value}
            className={`category-tab${categoria.value === cat.value ? ' active' : ''}`}
            onClick={() => setCategoria(cat)}
          >
            {cat.label}
          </button>
        ))}
      </div>

      {/* ── Grilla de eventos ─────────────────────────────── */}
      <section className="events-section">
        {loading && (
          <div className="state-box">
            <span className="spinner" />
            <p>Cargando eventos...</p>
          </div>
        )}

        {!loading && error && (
          <div className="state-box" style={{ color: 'var(--danger)' }}>
            <p>⚠️ {error}</p>
            <button className="btn-retry" onClick={() => fetchEventos(search, categoria.value)}>Reintentar</button>
          </div>
        )}

        {!loading && !error && eventos.length === 0 && (
          <div className="state-box">
            <p>🎭 No se encontraron eventos con esos filtros.</p>
            <button className="btn-retry" onClick={handleClearFilters}>Limpiar filtros</button>
          </div>
        )}

        {!loading && !error && eventos.length > 0 && (
          <>
            <p className="events-count">
              {eventos.length} evento{eventos.length !== 1 ? 's' : ''} encontrado{eventos.length !== 1 ? 's' : ''}
            </p>
            <div className="events-grid">
              {eventos.map((ev) => (
                <EventCard
                  key={ev.id}
                  evento={ev}
                  esFavorito={favoritoIds.has(ev.id)}
                  onToggleFavorito={toggleFavorito}
                  suscripcionActiva={suscripcion?.activa}
                />
              ))}
            </div>
          </>
        )}
      </section>

      {modalCancelarClub && (
        <ModalCancelarClub
          onClose={() => setModalCancelarClub(false)}
          onConfirm={handleCancelarSuscripcion}
          cancelando={cancelando}
        />
      )}
    </>
  )
}
