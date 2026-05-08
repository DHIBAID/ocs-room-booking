import { FormEvent, useEffect, useState } from 'react'
import { apiFetch } from '../api/client'
import type { Booking } from '../api/types'
import { useAuth } from '../hooks/useAuth'
import { useToast } from '../components/ToastProvider'

export const AdminPage = () => {
  const { token } = useAuth()
  const [bookings, setBookings] = useState<Booking[]>([])
  const { pushToast } = useToast()
  const [userForm, setUserForm] = useState({
    username: '',
    password: '',
    role: 'core',
  })
  const [roomForm, setRoomForm] = useState({
    block: '',
    name: '',
    capacity: '',
    status: 'available',
    allowed_purposes: '',
    notes: '',
  })

  const loadBookings = async () => {
    try {
      const data = await apiFetch('/api/admin/bookings', { token })
      setBookings(data)
    } catch (err) {
      pushToast('error', (err as Error).message)
    }
  }

  useEffect(() => {
    loadBookings()
  }, [])

  const handleCreateUser = async (event: FormEvent) => {
    event.preventDefault()
    try {
      await apiFetch('/api/admin/users', {
        method: 'POST',
        token,
        body: JSON.stringify(userForm),
      })
      pushToast('success', 'User account created.')
      setUserForm({ username: '', password: '', role: 'core' })
    } catch (err) {
      pushToast('error', (err as Error).message)
    }
  }

  const handleCreateRoom = async (event: FormEvent) => {
    event.preventDefault()
    try {
      await apiFetch('/api/admin/rooms', {
        method: 'POST',
        token,
        body: JSON.stringify({
          ...roomForm,
          capacity: Number(roomForm.capacity),
        }),
      })
      pushToast('success', 'Room added to inventory.')
      setRoomForm({
        block: '',
        name: '',
        capacity: '',
        status: 'available',
        allowed_purposes: '',
        notes: '',
      })
    } catch (err) {
      pushToast('error', (err as Error).message)
    }
  }

  return (
    <div className="page-body">
      <header className="page-header">
        <div>
          <span className="kicker">Admin control</span>
          <h2>Manage accounts, rooms, and bookings</h2>
          <p>Maintain access permissions and keep the schedule conflict-free.</p>
        </div>
      </header>

      <div className="page-grid">
        <section className="card">
          <h3>Create user</h3>
          <form className="form" onSubmit={handleCreateUser}>
            <label className="field">
              <span>Username</span>
              <input
                className="input"
                value={userForm.username}
                onChange={(event) =>
                  setUserForm({ ...userForm, username: event.target.value })
                }
                required
              />
            </label>
            <label className="field">
              <span>Password</span>
              <input
                className="input"
                type="password"
                value={userForm.password}
                onChange={(event) =>
                  setUserForm({ ...userForm, password: event.target.value })
                }
                required
              />
            </label>
            <label className="field">
              <span>Role</span>
              <select
                className="input"
                value={userForm.role}
                onChange={(event) =>
                  setUserForm({ ...userForm, role: event.target.value })
                }
              >
                <option value="core">Core</option>
                <option value="viewer">Viewer</option>
                <option value="admin">Admin</option>
              </select>
            </label>
            <button className="button" type="submit">
              Create user
            </button>
          </form>
        </section>

        <section className="card">
          <h3>Add room</h3>
          <form className="form" onSubmit={handleCreateRoom}>
            <label className="field">
              <span>Block</span>
              <input
                className="input"
                value={roomForm.block}
                onChange={(event) =>
                  setRoomForm({ ...roomForm, block: event.target.value })
                }
                required
              />
            </label>
            <label className="field">
              <span>Room name</span>
              <input
                className="input"
                value={roomForm.name}
                onChange={(event) =>
                  setRoomForm({ ...roomForm, name: event.target.value })
                }
                required
              />
            </label>
            <label className="field">
              <span>Capacity</span>
              <input
                className="input"
                type="number"
                min="1"
                value={roomForm.capacity}
                onChange={(event) =>
                  setRoomForm({ ...roomForm, capacity: event.target.value })
                }
                required
              />
            </label>
            <label className="field">
              <span>Allowed purposes</span>
              <input
                className="input"
                value={roomForm.allowed_purposes}
                onChange={(event) =>
                  setRoomForm({
                    ...roomForm,
                    allowed_purposes: event.target.value,
                  })
                }
                placeholder="OA,Interview,PPT"
              />
            </label>
            <label className="field">
              <span>Status</span>
              <select
                className="input"
                value={roomForm.status}
                onChange={(event) =>
                  setRoomForm({ ...roomForm, status: event.target.value })
                }
              >
                <option value="available">Available</option>
                <option value="unavailable">Unavailable</option>
              </select>
            </label>
            <label className="field">
              <span>Notes</span>
              <input
                className="input"
                value={roomForm.notes}
                onChange={(event) =>
                  setRoomForm({ ...roomForm, notes: event.target.value })
                }
                placeholder="Special constraints"
              />
            </label>
            <button className="button" type="submit">
              Add room
            </button>
          </form>
        </section>
      </div>

      <section className="card">
        <h3>All bookings</h3>
        <div className="table">
          {bookings.map((booking) => (
            <div key={booking.id} className="row">
              <div>
                <strong>{booking.room?.name ?? 'Room'}</strong>
                <span>{booking.room?.block}</span>
              </div>
              <div>
                <span>{booking.purpose}</span>
                <span>{booking.participants} participants</span>
              </div>
              <div>
                <span>{new Date(booking.start_time).toLocaleString()}</span>
                <span>to {new Date(booking.end_time).toLocaleString()}</span>
              </div>
              <div>
                <span className="pill status">{booking.status}</span>
              </div>
            </div>
          ))}
          {bookings.length === 0 && (
            <p className="empty">No bookings yet.</p>
          )}
        </div>
      </section>
    </div>
  )
}
