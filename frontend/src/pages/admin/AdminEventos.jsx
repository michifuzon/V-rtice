import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import client from '../../api/client'

const EMPTY_FORM = {
  titulo: '', descripcion: '', fecha_hora: '', duracion_minutos: '',
  ubicacion: '', categoria: '', cupo_total: '', precio: '', imagen_url: '',
}

export default function AdminEventos() {
  const [eventos, setEventos] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [feedback, setFeedback] = useState({ type: '', msg: '' })

  const [modal, setModal] = useState(false)
  const [editando, setEditando] = useState(null)
  const [form, setForm] = useState(EMPTY_FORM)
  const [saving, setSaving] = useState(false)
  const [formError, setFormError] = useState('')

  const [cancelando, setCancelando] = useState(null)
  const [deletingId, setDeletingId] = useState(null)

  useEffect(() => { fetchEventos() }, [])

  const fetchEventos = async () => {
    setLoading(true)
    setError('')
    try {
      const { data } = await client.get('/eventos')
      setEventos(Array.isArray(data) ? data : data.eventos ?? [])
    } catch {
      setError('No se pudieron cargar los eventos.')
    } finally {
      setLoading(false)
    }
  }

  const openCreate = () => {
    setEditando(null)
    setForm(EMPTY_FORM)
    setFormError('')
    setModal(true)
  }

  const openEdit = (ev) => {
    setEditando(ev)
    setForm({
      titulo: ev.titulo || ev.nombre || '',
      descripcion: ev.descripcion || '',
      fecha_hora: ev.fecha_hora ? ev.fecha_hora.slice(0, 16) : '',
      duracion_minutos: ev.duracion_minutos ?? '',
      ubicacion: ev.ubicacion || '',
      categoria: ev.categoria || '',
      cupo_total: ev.cupo_total ?? '',
      precio: ev.precio ?? '',
      imagen_url: ev.imagen_url || '',
    })
    setFormError('')
    setModal(true)
  }

  const closeModal = () => {
    setModal(false)
    setEditando(null)
    setForm(EMPTY_FORM)
    setFormError('')
  }

  const handleFormChange = (e) => {
    const { name, value } = e.target
    setForm(prev => ({ ...prev, [name]: value }))
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    setSaving(true)
    setFormError('')
    try {
      const payload = {
        titulo: form.titulo,
        descripcion: form.descripcion,
        fecha_hora: new Date(form.fecha_hora).toISOString(),
        duracion_minutos: Number(form.duracion_minutos) || 0,
        ubicacion: form.ubicacion,
        categoria: form.categoria,
        cupo_total: Number(form.cupo_total),
        precio: Number(form.precio) || 0,
        imagen_url: form.imagen_url,
      }
      if (editando) {
        await client.put(`/admin/eventos/${editando.id}`, payload)
        setFeedback({ type: 'success', msg: 'Evento actualizado correctamente.' })
      } else {
        await client.post('/admin/eventos', payload)
        setFeedback({ type: 'success', msg: 'Evento creado correctamente.' })
      }
      closeModal()
      fetchEventos()
    } catch (err) {
      setFormError(err.response?.data?.error || err.response?.data?.message || 'Error al guardar el evento.')
    } finally {
      setSaving(false)
    }
  }

  const handleCancelarEvento = async (id) => {
    setDeletingId(id)
    try {
      await client.delete(`/admin/eventos/${id}`)
      setFeedback({ type: 'success', msg: 'Evento cancelado.' })
      setEventos(prev => prev.filter(ev => ev.id !== id))
    } catch (err) {
      setFeedback({ type: 'error', msg: err.response?.data?.error || 'No se pudo cancelar el evento.' })
    } finally {
      setDeletingId(null)
      setCancelando(null)
    }
  }

  return (
    <div className="admin-page">
      <div className="admin-inner">
        <div className="admin-header">
          <div>
            <h1 className="admin-title">Panel de administración</h1>
            <p className="admin-subtitle">Gestión de eventos · {eventos.length} evento{eventos.length !== 1 ? 's' : ''}</p>
          </div>
          <button className="btn-pill btn-pill-solid" onClick={openCreate}>+ Nuevo evento</button>
        </div>

        {feedback.msg && (
          <div className={`alert alert-${feedback.type}`}>
            {feedback.msg}
            <button className="alert-close" onClick={() => setFeedback({ type: '', msg: '' })}>✕</button>
          </div>
        )}

        {loading && (
          <div className="state-box"><span className="spinner" /><p>Cargando eventos...</p></div>
        )}

        {!loading && error && (
          <div className="state-box" style={{ color: 'var(--danger)' }}>
            <p>⚠️ {error}</p>
            <button className="btn-retry" onClick={fetchEventos}>Reintentar</button>
          </div>
        )}

        {!loading && !error && (
          <div className="admin-table-wrap">
            <table className="admin-table">
              <thead>
                <tr>
                  <th>Título</th>
                  <th>Categoría</th>
                  <th>Fecha</th>
                  <th>Precio</th>
                  <th>Cupo total</th>
                  <th>Acciones</th>
                </tr>
              </thead>
              <tbody>
                {eventos.length === 0 && (
                  <tr>
                    <td colSpan="6" style={{ textAlign: 'center', padding: '2rem', color: 'var(--text-muted)' }}>
                      No hay eventos cargados.
                    </td>
                  </tr>
                )}
                {eventos.map(ev => (
                  <tr key={ev.id}>
                    <td className="admin-td-title">{ev.titulo || ev.nombre}</td>
                    <td>{ev.categoria || '—'}</td>
                    <td>
                      {ev.fecha_hora
                        ? new Date(ev.fecha_hora).toLocaleDateString('es-AR', { day: 'numeric', month: 'short', year: 'numeric' })
                        : '—'}
                    </td>
                    <td>
                      {ev.precio === 0 ? 'Gratis' : ev.precio ? `$${Number(ev.precio).toLocaleString('es-AR')}` : '—'}
                    </td>
                    <td>{ev.cupo_total ?? '—'}</td>
                    <td>
                      <div className="admin-actions">
                        <button className="btn-action btn-action-transfer" onClick={() => openEdit(ev)}>
                          Editar
                        </button>
                        <Link to={`/admin/reporte/${ev.id}`} className="btn-action btn-action-transfer">
                          Reporte
                        </Link>
                        {cancelando === ev.id ? (
                          <>
                            <span className="admin-confirm-label">¿Confirmar?</span>
                            <button
                              className="btn-action btn-action-cancel"
                              onClick={() => handleCancelarEvento(ev.id)}
                              disabled={deletingId === ev.id}
                            >
                              {deletingId === ev.id ? '...' : 'Sí'}
                            </button>
                            <button className="btn-action btn-action-transfer" onClick={() => setCancelando(null)}>
                              No
                            </button>
                          </>
                        ) : (
                          <button className="btn-action btn-action-cancel" onClick={() => setCancelando(ev.id)}>
                            Cancelar
                          </button>
                        )}
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {/* ── Modal crear / editar ──────────────────────────────── */}
      {modal && (
        <div className="modal-overlay" onClick={closeModal}>
          <div className="modal-box modal-box-lg" onClick={e => e.stopPropagation()}>
            <div className="modal-header">
              <h3>{editando ? 'Editar evento' : 'Nuevo evento'}</h3>
              <button className="modal-close-btn" onClick={closeModal}>✕</button>
            </div>
            <form onSubmit={handleSubmit}>
              <div className="modal-body modal-form-grid">
                {formError && (
                  <div className="form-error" style={{ gridColumn: '1/-1' }}>{formError}</div>
                )}

                <div className="field" style={{ gridColumn: '1/-1' }}>
                  <label>Título *</label>
                  <input name="titulo" value={form.titulo} onChange={handleFormChange} required placeholder="Nombre del evento" />
                </div>

                <div className="field" style={{ gridColumn: '1/-1' }}>
                  <label>Descripción</label>
                  <textarea
                    name="descripcion"
                    value={form.descripcion}
                    onChange={handleFormChange}
                    className="admin-textarea"
                    placeholder="Descripción del evento..."
                  />
                </div>

                <div className="field">
                  <label>Fecha y hora *</label>
                  <input type="datetime-local" name="fecha_hora" value={form.fecha_hora} onChange={handleFormChange} required />
                </div>

                <div className="field">
                  <label>Duración (minutos)</label>
                  <input type="number" name="duracion_minutos" value={form.duracion_minutos} onChange={handleFormChange} min="1" placeholder="120" />
                </div>

                <div className="field">
                  <label>Ubicación</label>
                  <input name="ubicacion" value={form.ubicacion} onChange={handleFormChange} placeholder="Teatro Gran Rex, CABA" />
                </div>

                <div className="field">
                  <label>Categoría</label>
                  <input name="categoria" value={form.categoria} onChange={handleFormChange} placeholder="música, teatro..." />
                </div>

                <div className="field">
                  <label>Cupo total *</label>
                  <input type="number" name="cupo_total" value={form.cupo_total} onChange={handleFormChange} min="1" required placeholder="500" />
                </div>

                <div className="field">
                  <label>Precio</label>
                  <input type="number" name="precio" value={form.precio} onChange={handleFormChange} min="0" placeholder="0 = Gratis" />
                </div>

                <div className="field" style={{ gridColumn: '1/-1' }}>
                  <label>URL de imagen</label>
                  <input name="imagen_url" value={form.imagen_url} onChange={handleFormChange} placeholder="https://..." />
                </div>
              </div>

              <div className="modal-footer">
                <button type="button" className="btn-modal-cancel" onClick={closeModal}>Cancelar</button>
                <button type="submit" className="btn-modal-confirm" disabled={saving}>
                  {saving ? 'Guardando...' : editando ? 'Guardar cambios' : 'Crear evento'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}
