import { useState, useEffect } from 'react'
import { useParams, useNavigate, Link } from 'react-router-dom'
import client from '../api/client'

const EMOJI = {
  'Música': '🎵', 'Humor': '😂', 'Teatro': '🎭', 'Deportes': '⚽',
  'Arte': '🎨', 'Cine': '🎬', 'Tecnología': '💻',
}
const getEmoji = (cat) => EMOJI[cat] ?? '🎉'

export default function EventoDetalle() {
  const { id }   = useParams()
  const navigate = useNavigate()
  const token    = localStorage.getItem('token')

  const [evento,   setEvento]   = useState(null)
  const [loading,  setLoading]  = useState(true)
  const [error,    setError]    = useState('')

  const [cantidad,    setCantidad]    = useState(1)
  const [comprando,   setComprando]   = useState(false)
  const [compraOk,    setCompraOk]    = useState(false)
  const [compraError, setCompraError] = useState('')

  const [esFavorito,   setEsFavorito]   = useState(false)
  const [togglingFav,  setTogglingFav]  = useState(false)
  const [suscripcion,  setSuscripcion]  = useState(null)

  useEffect(() => {
    const fetchEvento = async () => {
      setLoading(true)
      setError('')
      try {
        const { data } = await client.get(`/eventos/${id}`)
        setEvento(data.evento ?? data)
      } catch (err) {
        setError(err.response?.status === 404 ? 'El evento no existe.' : 'No se pudo cargar el evento.')
      } finally {
        setLoading(false)
      }
    }
    fetchEvento()
  }, [id])

  useEffect(() => {
    if (!token) return
    client.get('/favoritos')
      .then(({ data }) => {
        const ids = (data.favoritos ?? []).map(fav => fav.evento_id)
        setEsFavorito(ids.includes(Number(id)))
      })
      .catch(() => {})
    client.get('/club/estado')
      .then(({ data }) => setSuscripcion(data))
      .catch(() => {})
  }, [id, token])

  const handleToggleFavorito = async () => {
    if (!token) { navigate('/login'); return }
    setTogglingFav(true)
    try {
      if (esFavorito) {
        await client.delete(`/favoritos/${id}`)
        setEsFavorito(false)
      } else {
        await client.post(`/favoritos/${id}`)
        setEsFavorito(true)
      }
    } catch { /* ignore */ } finally {
      setTogglingFav(false)
    }
  }

  const handleComprar = async () => {
    if (!token) { navigate('/login'); return }
    setComprando(true)
    setCompraError('')
    try {
      await client.post('/entradas', { evento_id: Number(id), cantidad })
      setCompraOk(true)
    } catch (err) {
      setCompraError(
        err.response?.data?.error ||
        err.response?.data?.message ||
        'No se pudo completar la compra.'
      )
    } finally {
      setComprando(false)
    }
  }

  if (loading) return (
    <div className="detalle-page">
      <div className="state-box"><span className="spinner" /><p>Cargando evento...</p></div>
    </div>
  )

  if (error) return (
    <div className="detalle-page">
      <div className="state-box" style={{ color: 'var(--danger)' }}>
        <p>⚠️ {error}</p>
        <Link to="/" className="btn-retry">← Volver al catálogo</Link>
      </div>
    </div>
  )

  const nombre = evento.nombre || evento.titulo || 'Sin nombre'
  const emoji  = getEmoji(evento.categoria)

  const fechaFmt = evento.fecha_hora
    ? new Date(evento.fecha_hora).toLocaleDateString('es-AR', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' })
    : null
  const horaFmt = evento.fecha_hora
    ? new Date(evento.fecha_hora).toLocaleTimeString('es-AR', { hour: '2-digit', minute: '2-digit' })
    : null

  const precioBase    = evento.precio ?? 0
  const hayDescuento  = suscripcion?.activa && precioBase > 0
  const precioUnitario = hayDescuento ? Math.round(precioBase * 0.9 * 100) / 100 : precioBase
  const precioTotal   = precioUnitario * cantidad

  const precioLabel = precioBase === 0 ? 'Gratis' : `$${Number(precioBase).toLocaleString('es-AR')}`
  const totalLabel  = precioBase === 0 ? 'Gratis' : `$${Number(precioTotal).toLocaleString('es-AR')}`

  return (
    <div className="detalle-page">
      <div className="detalle-inner">
        <Link to="/" className="back-link">← Volver al catálogo</Link>

        <div className="detalle-grid">
          {/* ── Columna principal ─────────────────────────── */}
          <div className="detalle-main">
            <div className="detalle-placeholder">{emoji}</div>

            <div className="detalle-content">
              {evento.categoria && (
                <div className="event-badges">
                  {evento.categoria.split(',').map(c => c.trim()).map(c => (
                    <span key={c} className="detalle-badge">{c}</span>
                  ))}
                </div>
              )}

              <h1 className="detalle-title">{nombre}</h1>

              <div className="detalle-meta">
                {fechaFmt && (
                  <div className="detalle-meta-item">
                    <span>📅</span>
                    <span>{fechaFmt}{horaFmt ? ` · ${horaFmt} hs` : ''}</span>
                  </div>
                )}
                {evento.ubicacion && (
                  <div className="detalle-meta-item">
                    <span>📍</span>
                    <span>{evento.ubicacion}</span>
                  </div>
                )}
                {evento.cupo_disponible != null && (
                  <div className="detalle-meta-item">
                    <span>👥</span>
                    <span>Lugares disponibles: {evento.cupo_disponible}</span>
                  </div>
                )}
              </div>

              {evento.descripcion && (
                <div className="detalle-description">
                  <h3>Descripción</h3>
                  <p>{evento.descripcion}</p>
                </div>
              )}
            </div>
          </div>

          {/* ── Sidebar de compra ─────────────────────────── */}
          <aside className="detalle-sidebar">

            {/* Precio */}
            <div className="sidebar-precio-wrap">
              {hayDescuento ? (
                <>
                  <div className="sidebar-precio-descuento">
                    <span className="precio-tachado-lg">{precioLabel}</span>
                    <span className="badge-descuento">10% OFF</span>
                  </div>
                  <p className="sidebar-price">${Number(precioUnitario).toLocaleString('es-AR')}</p>
                  <p className="sidebar-precio-club">💜 Precio Club Vórtice</p>
                </>
              ) : (
                <>
                  <p className={`sidebar-price ${precioBase === 0 ? 'sidebar-price-free' : ''}`}>
                    {precioLabel}
                  </p>
                  {precioBase > 0 && <p className="sidebar-precio-label">por entrada</p>}
                </>
              )}
            </div>

            <hr className="sidebar-divider" />

            {/* Selector de cantidad */}
            {!compraOk && token && precioBase > 0 && (
              <div className="cantidad-wrap">
                <label className="cantidad-label">Cantidad de entradas</label>
                <div className="cantidad-control">
                  <button
                    className="cantidad-btn"
                    onClick={() => setCantidad(c => Math.max(1, c - 1))}
                    disabled={cantidad <= 1}
                  >−</button>
                  <span className="cantidad-valor">{cantidad}</span>
                  <button
                    className="cantidad-btn"
                    onClick={() => setCantidad(c => Math.min(10, c + 1))}
                    disabled={cantidad >= 10}
                  >+</button>
                </div>
                {cantidad > 1 && (
                  <p className="cantidad-total">
                    Total: <strong>{totalLabel}</strong>
                  </p>
                )}
              </div>
            )}

            {/* Feedback de compra */}
            {compraOk && (
              <div className="alert alert-success">
                ✅ {cantidad > 1 ? `${cantidad} entradas compradas` : '¡Entrada comprada!'}
                {' '}<Link to="/mis-entradas">Ver mis entradas →</Link>
              </div>
            )}

            {compraError && (
              <div className="alert alert-error">⚠️ {compraError}</div>
            )}

            {/* Botón comprar / login */}
            {!compraOk && (
              token ? (
                <button className="sidebar-btn" onClick={handleComprar} disabled={comprando}>
                  {comprando ? 'Procesando...' : cantidad > 1 ? `🎟️ Comprar ${cantidad} entradas` : '🎟️ Comprar entrada'}
                </button>
              ) : (
                <>
                  <p className="sidebar-login-msg">
                    <Link to="/login">Iniciá sesión</Link> para comprar entradas.
                  </p>
                  <Link to="/login" className="sidebar-btn" style={{ textAlign: 'center' }}>
                    Iniciar sesión
                  </Link>
                </>
              )
            )}

            {/* Favorito */}
            {token && (
              <button
                className={`btn-favorito${esFavorito ? ' btn-favorito-active' : ''}`}
                onClick={handleToggleFavorito}
                disabled={togglingFav}
              >
                {togglingFav ? '...' : esFavorito ? '♥ En favoritos' : '♡ Agregar a favoritos'}
              </button>
            )}

          </aside>
        </div>
      </div>
    </div>
  )
}
