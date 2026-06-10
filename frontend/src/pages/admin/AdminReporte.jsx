import { useState, useEffect } from 'react'
import { useParams, Link } from 'react-router-dom'
import client from '../../api/client'

export default function AdminReporte() {
  const { id } = useParams()
  const [reporte, setReporte] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    const fetchReporte = async () => {
      setLoading(true)
      try {
        const { data } = await client.get(`/admin/eventos/${id}/reporte`)
        setReporte(data)
      } catch (err) {
        setError(
          err.response?.status === 403
            ? 'No tenés permisos para ver este reporte.'
            : 'No se pudo cargar el reporte.'
        )
      } finally {
        setLoading(false)
      }
    }
    fetchReporte()
  }, [id])

  if (loading) return (
    <div className="admin-page">
      <div className="state-box"><span className="spinner" /><p>Cargando reporte...</p></div>
    </div>
  )

  if (error) return (
    <div className="admin-page">
      <div className="state-box" style={{ color: 'var(--danger)' }}>
        <p>⚠️ {error}</p>
        <Link to="/admin" className="btn-retry">← Volver al panel</Link>
      </div>
    </div>
  )

  const pct = reporte.porcentaje_ocupacion ?? 0

  return (
    <div className="admin-page">
      <div className="admin-inner">
        <Link to="/admin" className="back-link">← Volver al panel</Link>

        <div className="admin-header">
          <div>
            <h1 className="admin-title">{reporte.titulo}</h1>
            <p className="admin-subtitle">Reporte de ocupación · Evento #{reporte.evento_id}</p>
          </div>
        </div>

        {/* Tarjetas de estadísticas */}
        <div className="reporte-stats">
          <div className="reporte-stat-card">
            <span className="reporte-stat-value">{reporte.cupo_total}</span>
            <span className="reporte-stat-label">Cupo total</span>
          </div>
          <div className="reporte-stat-card">
            <span className="reporte-stat-value">{reporte.entradas_vendidas}</span>
            <span className="reporte-stat-label">Entradas vendidas</span>
          </div>
          <div className="reporte-stat-card">
            <span className="reporte-stat-value">{reporte.cupo_disponible}</span>
            <span className="reporte-stat-label">Cupo disponible</span>
          </div>
          <div className="reporte-stat-card reporte-stat-highlight">
            <span className="reporte-stat-value">{pct.toFixed(1)}%</span>
            <span className="reporte-stat-label">Ocupación</span>
          </div>
        </div>

        {/* Barra de progreso */}
        <div className="reporte-bar-wrap">
          <div className="reporte-bar">
            <div className="reporte-bar-fill" style={{ width: `${Math.min(pct, 100)}%` }} />
          </div>
          <span className="reporte-bar-label">{pct.toFixed(1)}% vendido</span>
        </div>

        {/* Tabla de compradores */}
        <h2 className="section-heading" style={{ marginTop: '2rem' }}>
          Compradores ({reporte.compradores?.length ?? 0})
        </h2>

        {(!reporte.compradores || reporte.compradores.length === 0) ? (
          <p style={{ color: 'var(--text-muted)', fontSize: '.9rem' }}>Aún no hay compradores para este evento.</p>
        ) : (
          <div className="admin-table-wrap">
            <table className="admin-table">
              <thead>
                <tr>
                  <th>Nombre</th>
                  <th>Apellido</th>
                  <th>Email</th>
                  <th>Fecha de compra</th>
                </tr>
              </thead>
              <tbody>
                {reporte.compradores.map((c, i) => (
                  <tr key={i}>
                    <td>{c.nombre}</td>
                    <td>{c.apellido}</td>
                    <td>{c.email}</td>
                    <td>
                      {new Date(c.fecha_compra).toLocaleString('es-AR', {
                        day: 'numeric', month: 'short', year: 'numeric',
                        hour: '2-digit', minute: '2-digit',
                      })}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  )
}
