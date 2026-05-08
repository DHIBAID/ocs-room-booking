import { NavLink, Outlet } from 'react-router-dom'
import { useAuth } from '../hooks/useAuth'

export const AppShell = () => {
  const { user, logout } = useAuth()

  return (
    <div className="app-shell">
      <aside className="side-panel">
        <div className="brand-block">
          <span className="kicker">OCS IITH</span>
          <h1>Room Booking</h1>
          <p>Coordinate OA, interviews, and PPTs with clarity.</p>
        </div>

        <nav className="nav">
          <NavLink to="/dashboard" className="nav-link">
            Availability
          </NavLink>
          <NavLink to="/bookings" className="nav-link">
            My bookings
          </NavLink>
          {user?.role === 'admin' && (
            <NavLink to="/admin" className="nav-link">
              Admin control
            </NavLink>
          )}
        </nav>

        <div className="profile">
          <div>
            <span className="pill">{user?.role ?? 'user'}</span>
            <strong>{user?.username}</strong>
          </div>
          <button type="button" className="button ghost" onClick={logout}>
            Log out
          </button>
        </div>
      </aside>

      <section className="page">
        <Outlet />
      </section>
    </div>
  )
}
