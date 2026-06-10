import { useState, useEffect, useCallback } from 'react'
import { Link } from 'react-router-dom'
import client from '../api/client'

const EMOJI = {
  'música': '🎵', 'humor': '😂', 'teatro': '🎭', 'deportes': '⚽',
  'arte': '🎨', 'cine': '🎬', 'tecnología': '💻', 'espectáculo': '🎪',
}
const getEmoji = (cat) => {
  if (!cat) return '🎉'
  const key = cat.split(',')[0].trim().toLowerCase()
  return EMOJI[key] ?? '🎉'
}

/* ── Modal de cancelación ────────────────────────────────── */
function ModalCancelar({ entrada, onClose, onConfirm }) {
  const [loading, setLoading] = useState(false)
  const evento = entrada.evento || {}
  const nombre = evento.titulo || evento.nombre || `Entrada #${entrada.id}`

  const handleConfirm = async () => {
    setLoading(true)
    await onConfirm(entrada.id)
    setLoading(false)
  }

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-box" onClick={e => e.stopPropagation()}>
        <div className="modal-header">
          <h3>Cancelar entrada</h3>
          <button className="modal-close-btn" onClick={onClose}>✕</button>
        </div>
        <div className="modal-body">
          <p className="modal-desc">
            ¿Confirmás la cancelación de tu entrada para <strong>{nombre}</strong>?
            Esta acción no se puede deshacer.
          </p>
          <div className="modal-footer">
            <button className="btn-modal-cancel" onClick={onClose}>No, volver</button>
            <button className="btn-modal-confirm btn-modal-danger" onClick={handleConfirm} disabled={loading}>
              {loading ? 'Cancelando...' : 'Sí, cancelar entrada'}
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}

/* ── Modal de traspaso ───────────────────────────────────── */
function ModalTraspaso({ entradaId, onClose, onSuccess }) {
  const [email, setEmail] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const handleSubmit = async (e) => {
    e.preventDefault()
    setLoading(true)
    setError('')
    try {
      await client.post(`/entradas/${entradaId}/transfer`, { email })
      onSuccess()
    } catch (err) {
      setError(
        err.response?.data?.error ||
        err.response?.data?.message ||
        'No se pudo realizar el traspaso.'
      )
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-box" onClick={e => e.stopPropagation()}>
        <div className="modal-header">
          <h3>Traspasar entrada</h3>
          <button className="modal-close-btn" onClick={onClose}>✕</button>
        </div>
        <form onSubmit={handleSubmit}>
          <div className="modal-body">
            <p className="modal-desc">
              Ingresá el email del usuario al que querés transferir esta entrada.
              La entrada dejará de ser tuya una vez confirmado.
            </p>
            <div className="field">
              <label htmlFor="email-destino">Email del destinatario</label>
              <input
                id="email-destino"
                type="email"
                value={email}
                onChange={e => { setEmail(e.target.value); setError('') }}
                placeholder="destinatario@email.com"
                required
                autoFocus
              />
            </div>
            {error && <p className="form-error">{error}</p>}
            <div className="modal-footer">
              <button type="button" className="btn-modal-cancel" onClick={onClose}>Cancelar</button>
              <button type="submit" className="btn-modal-confirm" disabled={loading}>
                {loading ? 'Transfiriendo...' : 'Confirmar traspaso'}
              </button>
            </div>
          </div>
        </form>
      </div>
    </div>
  )
}

/* ── EntradaCard ─────────────────────────────────────────── */
function EntradaCard({ entrada, onCancelar, onTraspasar }) {
  const evento = entrada.evento || {}
  const nombre = evento.titulo || evento.nombre || `Evento #${entrada.evento_id}`
  const emoji = getEmoji(evento.categoria)

  const fechaEvento = evento.fecha_hora
    ? new Date(evento.fecha_hora).toLocaleDateString('es-AR', {
        weekday: 'short', day: 'numeric', month: 'short', year: 'numeric',
      })
    : null

  const fechaCompra = entrada.fecha_compra || entrada.created_at
  const fechaCompraFmt = fechaCompra
    ? new Date(fechaCompra).toLocaleDateString('es-AR')
    : null

  const precioLabel = entrada.precio_pagado != null
    ? entrada.precio_pagado === 0 ? 'Gratis' : `$${Number(entrada.precio_pagado).toLocaleString('es-AR')}`
    : null

  const esCancelada = entrada.estado === 'cancelada'
  const esTraspasada = entrada.estado === 'traspasada'
  const esActiva = !esCancelada && !esTraspasada

  const badgeClass = esCancelada ? 'badge-cancelada' : esTraspasada ? 'badge-traspasada' : 'badge-activa'
  const badgeLabel = esCancelada ? 'Cancelada' : esTraspasada ? 'Traspasada' : 'Activa'

  return (
    <div className={`entrada-card${!esActiva ? ' cancelada' : ''}`}>
      <div className="entrada-icon-wrap">
        <span className="entrada-icon">{emoji}</span>
        <span className={`entrada-status-badge ${badgeClass}`}>{badgeLabel}</span>
      </div>

      <div className="entrada-info">
        <p className="entrada-nombre">
          <Link to={`/eventos/${entrada.evento_id}`}>{nombre}</Link>
        </p>
        {fechaEvento && <p className="entrada-meta">📅 {fechaEvento}</p>}
        {evento.ubicacion && <p className="entrada-meta">📍 {evento.ubicacion}</p>}
        {precioLabel && <p className="entrada-meta" style={{ color: 'var(--accent-dark)', fontWeight: 600 }}>💳 {precioLabel}</p>}
        {fechaCompraFmt && <p className="entrada-compra-date">Comprada el {fechaCompraFmt}</p>}
      </div>

      {esActiva && (
        <div className="entrada-actions">
          <button className="btn-action btn-action-transfer" onClick={() => onTraspasar(entrada.id)}>
            ↔️ Traspasar
          </button>
          <button className="btn-action btn-action-cancel" onClick={() => onCancelar(entrada)}>
            ✕ Cancelar
          </button>
        </div>
      )}
    </div>
  )
}

/* ── MisEntradas ─────────────────────────────────────────── */
export default function MisEntradas() {
  const [entradas, setEntradas] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [traspasarId, setTraspasarId] = useState(null)
  const [cancelarEntrada, setCancelarEntrada] = useState(null)
  const [feedback, setFeedback] = useState(null)

  const fetchEntradas = useCallback(async () => {
    setLoading(true)
    setError('')
    try {
      const { data } = await client.get('/entradas/me')
      const lista = Array.isArray(data) ? data : data.entradas ?? []
      setEntradas(lista)
    } catch {
      setError('No se pudieron cargar tus entradas.')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => { fetchEntradas() }, [fetchEntradas])

  const handleCancelar = async (entradaId) => {
    try {
      await client.delete(`/entradas/${entradaId}`)
      setCancelarEntrada(null)
      setFeedback({ type: 'success', msg: '✅ Entrada cancelada. El cupo fue devuelto al evento.' })
      fetchEntradas()
    } catch (err) {
      setCancelarEntrada(null)
      setFeedback({
        type: 'error',
        msg: err.response?.data?.error || err.response?.data?.message || 'No se pudo cancelar la entrada.',
      })
    }
  }

  const handleTraspasoOk = () => {
    setTraspasarId(null)
    setFeedback({ type: 'success', msg: '✅ Entrada traspasada correctamente.' })
    fetchEntradas()
  }

  const activas = entradas.filter(e => e.estado !== 'cancelada' && e.estado !== 'traspasada')
  const inactivas = entradas.filter(e => e.estado === 'cancelada' || e.estado === 'traspasada')

  return (
    <div className="mis-entradas-page">
      <div className="mis-entradas-inner">
        <div className="page-title-row">
          <h1>Mis Entradas</h1>
          <Link to="/" className="btn-pill btn-pill-outline" style={{ fontSize: '.85rem', padding: '7px 16px' }}>
            Ver eventos
          </Link>
        </div>

        {feedback && (
          <div className={`alert alert-${feedback.type}`}>
            {feedback.msg}
            <button className="alert-close" onClick={() => setFeedback(null)}>✕</button>
          </div>
        )}

        {loading && <div className="state-box"><span className="spinner" /><p>Cargando tus entradas...</p></div>}

        {!loading && error && (
          <div className="state-box" style={{ color: 'var(--danger)' }}>
            <p>⚠️ {error}</p>
            <button className="btn-retry" onClick={fetchEntradas}>Reintentar</button>
          </div>
        )}

        {!loading && !error && entradas.length === 0 && (
          <div className="state-box">
            <p>🎟️ Todavía no compraste ninguna entrada.</p>
            <Link to="/" className="btn-pill btn-pill-solid" style={{ fontSize: '.875rem' }}>
              Ver eventos disponibles
            </Link>
          </div>
        )}

        {!loading && !error && activas.length > 0 && (
          <section>
            <p className="section-heading">Activas ({activas.length})</p>
            <div className="entradas-list">
              {activas.map(e => (
                <EntradaCard
                  key={e.id}
                  entrada={e}
                  onCancelar={setCancelarEntrada}
                  onTraspasar={setTraspasarId}
                />
              ))}
            </div>
          </section>
        )}

        {!loading && !error && inactivas.length > 0 && (
          <section>
            <p className="section-heading muted">Historial ({inactivas.length})</p>
            <div className="entradas-list">
              {inactivas.map(e => (
                <EntradaCard
                  key={e.id}
                  entrada={e}
                  onCancelar={setCancelarEntrada}
                  onTraspasar={setTraspasarId}
                />
              ))}
            </div>
          </section>
        )}
      </div>

      {cancelarEntrada && (
        <ModalCancelar
          entrada={cancelarEntrada}
          onClose={() => setCancelarEntrada(null)}
          onConfirm={handleCancelar}
        />
      )}

      {traspasarId && (
        <ModalTraspaso
          entradaId={traspasarId}
          onClose={() => setTraspasarId(null)}
          onSuccess={handleTraspasoOk}
        />
      )}
    </div>
  )
}
