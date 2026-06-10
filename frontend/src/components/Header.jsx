import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import vortice_icon from '../assets/vortice_favicon.svg'

export default function Header() {
  const navigate = useNavigate()
  const [open, setOpen] = useState(false)
  const token = localStorage.getItem('token')
  const usuario = JSON.parse(localStorage.getItem('usuario') || 'null')

  const handleLogout = () => {
    localStorage.removeItem('token')
    localStorage.removeItem('usuario')
    navigate('/login')
    setOpen(false)
  }

  const close = () => setOpen(false)

  return (
    <header className="header">
      <Link to="/" className="header-logo" onClick={close}>
        <img src={vortice_icon} alt="" className="header-logo-img" />
        <span className="header-logo-text">Vórtice</span>
      </Link>

      {/* Nav desktop */}
      <nav className="header-nav">
        <Link to="/" className="header-nav-link">Eventos</Link>
        {token && <Link to="/mis-entradas" className="header-nav-link">Mis entradas</Link>}
        {token && <Link to="/mis-favoritos" className="header-nav-link">Mis favoritos</Link>}
        {token && usuario?.rol === 'admin' && (
          <Link to="/admin" className="btn-pill btn-pill-outline header-nav-admin">Panel Admin</Link>
        )}
        {token ? (
          <>
            <span className="header-user">Hola, {usuario?.nombre ?? 'Usuario'}</span>
            <button className="btn-pill btn-pill-outline" onClick={handleLogout}>
              Cerrar sesión
            </button>
          </>
        ) : (
          <>
            <Link to="/login" className="btn-pill btn-pill-outline">Iniciar sesión</Link>
            <Link to="/register" className="btn-pill btn-pill-solid">Registrarse</Link>
          </>
        )}
      </nav>

      {/* Botón hamburguesa (solo mobile) */}
      <button
        className={`hamburger${open ? ' is-open' : ''}`}
        onClick={() => setOpen(!open)}
        aria-label="Menú"
      >
        <span /><span /><span />
      </button>

      {/* Menú mobile */}
      {open && (
        <div className="mobile-menu">
          <Link to="/" className="mobile-nav-link" onClick={close}>Eventos</Link>
          {token && <Link to="/mis-entradas" className="mobile-nav-link" onClick={close}>Mis entradas</Link>}
          {token && <Link to="/mis-favoritos" className="mobile-nav-link" onClick={close}>Mis favoritos</Link>}
          {token && usuario?.rol === 'admin' && (
            <Link to="/admin" className="mobile-nav-link" onClick={close} style={{ color: '#fbbf24' }}>Panel Admin</Link>
          )}
          {token ? (
            <>
              <span className="mobile-user">Hola, {usuario?.nombre ?? 'Usuario'}</span>
              <button className="btn-pill btn-pill-outline mobile-btn" onClick={handleLogout}>
                Cerrar sesión
              </button>
            </>
          ) : (
            <div className="mobile-actions">
              <Link to="/login" className="btn-pill btn-pill-outline mobile-btn" onClick={close}>
                Iniciar sesión
              </Link>
              <Link to="/register" className="btn-pill btn-pill-solid mobile-btn" onClick={close}>
                Registrarse
              </Link>
            </div>
          )}
        </div>
      )}
    </header>
  )
}
