import { useState, useEffect } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import client from '../api/client'

const EMOJI = {
  'música': '🎵', 'humor': '😂', 'teatro': '🎭', 'deportes': '⚽',
  'arte': '🎨', 'cine': '🎬', 'tecnología': '💻', 'espectáculo': '🎪',
}
const getEmoji = (cat) => {
  if (!cat) return '🎉'
  const cats = cat.split(',').map(c => c.trim().toLowerCase())
  for (const c of cats) { if (EMOJI[c]) return EMOJI[c] }
  return '🎉'
}

const primerCategoria = (cat) => {
  if (!cat) return 'Sin categoría'
  return cat.split(',')[0].trim() || 'Sin categoría'
}

const cap = (s) => s ? s.charAt(0).toUpperCase() + s.slice(1) : s

const ORDEN_CATS = ['música', 'teatro', 'humor', 'deportes', 'arte', 'espectáculo', 'cine', 'tecnología']

function agruparPorCategoria(favoritos) {
  const grupos = {}
  for (const fav of favoritos) {
    const cat = primerCategoria(fav.evento?.categoria)
    if (!grupos[cat]) grupos[cat] = []
    grupos[cat].push(fav)
  }
  return Object.entries(grupos).sort(([a], [b]) => {
    const ia = ORDEN_CATS.indexOf(a.toLowerCase())
    const ib = ORDEN_CATS.indexOf(b.toLowerCase())
    if (ia === -1 && ib === -1) return a.localeCompare(b)
    if (ia === -1) return 1
    if (ib === -1) return -1
    return ia - ib
  })
}

function FavCard({ fav, removing, onRemove }) {
  const ev = fav.evento || {}
  const nombre = ev.titulo || ev.nombre || 'Sin nombre'
  const emoji = getEmoji(ev.categoria)
  const fecha = ev.fecha_hora
    ? new Date(ev.fecha_hora).toLocaleDateString('es-AR', { day: 'numeric', month: 'short', year: 'numeric' })
    : null
  const precioLabel = ev.precio == null ? null : ev.precio === 0 ? 'Gratis' : `$${Number(ev.precio).toLocaleString('es-AR')}`

  return (
    <div className="event-card-wrap">
      <Link to={`/eventos/${fav.evento_id}`} className="event-card">
        <div className="event-placeholder">{emoji}</div>
        <div className="event-body">
          {ev.categoria && (
            <div className="event-badges">
              {ev.categoria.split(',').map(c => c.trim()).map(c => (
                <span key={c} className="event-category-badge">{c}</span>
              ))}
            </div>
          )}
          <h3 className="event-title">{nombre}</h3>
          {ev.ubicacion && <p className="event-venue">📍 {ev.ubicacion}</p>}
          {fecha && <p className="event-date">📅 {fecha}</p>}
          {precioLabel && (
            <p style={{ fontWeight: 700, color: 'var(--accent-dark)', fontSize: '.92rem', marginTop: '.2rem' }}>
              {precioLabel}
            </p>
          )}
        </div>
        <div className="event-footer">
          <span className="btn-buy">Ver evento</span>
        </div>
      </Link>
      <button
        className="heart-btn heart-active"
        onClick={() => onRemove(fav.evento_id)}
        disabled={removing === fav.evento_id}
        title="Quitar de favoritos"
      >
        {removing === fav.evento_id ? '…' : '♥'}
      </button>
    </div>
  )
}

export default function MisFavoritos() {
  const navigate = useNavigate()
  const token = localStorage.getItem('token')
  const [favoritos, setFavoritos] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [removing, setRemoving] = useState(null)
  const [feedback, setFeedback] = useState('')
  const [filtrocat, setFiltrocat] = useState('Todos')

  useEffect(() => {
    if (!token) { navigate('/login'); return }
    fetchFavoritos()
  }, [])

  const fetchFavoritos = async () => {
    setLoading(true)
    try {
      const { data } = await client.get('/favoritos')
      setFavoritos(data.favoritos ?? [])
    } catch {
      setError('No se pudieron cargar tus favoritos.')
    } finally {
      setLoading(false)
    }
  }

  const handleRemove = async (eventoId) => {
    setRemoving(eventoId)
    try {
      await client.delete(`/favoritos/${eventoId}`)
      setFavoritos(prev => prev.filter(f => f.evento_id !== eventoId))
      setFeedback('Eliminado de favoritos.')
      setTimeout(() => setFeedback(''), 3000)
    } catch {
      setFeedback('No se pudo eliminar el favorito.')
    } finally {
      setRemoving(null)
    }
  }

  const grupos = agruparPorCategoria(favoritos)
  const todasCats = ['Todos', ...grupos.map(([cat]) => cat)]

  const gruposFiltrados = filtrocat === 'Todos'
    ? grupos
    : grupos.filter(([cat]) => cat === filtrocat)

  return (
    <div className="mis-entradas-page">
      <div className="mis-entradas-inner" style={{ maxWidth: '1100px' }}>
        <div className="page-title-row">
          <h1>Mis favoritos ♥</h1>
          <Link to="/" className="back-link">← Ver catálogo</Link>
        </div>

        {feedback && (
          <div className="alert alert-success" style={{ marginBottom: '1rem' }}>
            {feedback}
            <button className="alert-close" onClick={() => setFeedback('')}>✕</button>
          </div>
        )}

        {loading && <div className="state-box"><span className="spinner" /><p>Cargando favoritos...</p></div>}

        {!loading && error && (
          <div className="state-box" style={{ color: 'var(--danger)' }}>
            <p>⚠️ {error}</p>
            <button className="btn-retry" onClick={fetchFavoritos}>Reintentar</button>
          </div>
        )}

        {!loading && !error && favoritos.length === 0 && (
          <div className="state-box">
            <p>No tenés eventos guardados todavía.</p>
            <Link to="/" className="btn-retry">Explorar eventos →</Link>
          </div>
        )}

        {!loading && !error && favoritos.length > 0 && (
          <>
            {/* Filtro por categoría */}
            {todasCats.length > 2 && (
              <div className="category-bar" style={{ marginBottom: '1.5rem', borderRadius: 'var(--radius)', border: '1px solid var(--border)' }}>
                {todasCats.map(cat => (
                  <button
                    key={cat}
                    className={`category-tab${filtrocat === cat ? ' active' : ''}`}
                    onClick={() => setFiltrocat(cat)}
                  >
                    {cat === 'Todos' ? `Todos (${favoritos.length})` : `${cap(cat)} (${grupos.find(([c]) => c === cat)?.[1].length ?? 0})`}
                  </button>
                ))}
              </div>
            )}

            {/* Grupos por categoría */}
            {gruposFiltrados.map(([cat, favs]) => (
              <div key={cat} className="favs-categoria-section">
                <h2 className="favs-categoria-titulo">
                  {getEmoji(cat)} {cap(cat)}
                  <span className="favs-categoria-count">{favs.length} evento{favs.length !== 1 ? 's' : ''}</span>
                </h2>
                <div className="events-grid">
                  {favs.map(fav => (
                    <FavCard
                      key={fav.evento_id}
                      fav={fav}
                      removing={removing}
                      onRemove={handleRemove}
                    />
                  ))}
                </div>
              </div>
            ))}
          </>
        )}
      </div>
    </div>
  )
}
