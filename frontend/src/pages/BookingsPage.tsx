import { useEffect, useState } from 'react'
import { apiFetch } from '../api/client'
import type { Booking } from '../api/types'
import { useAuth } from '../hooks/useAuth'
import { useToast } from '../components/ToastProvider'

export const BookingsPage = () => {
  const { token, user } = useAuth()
  const [bookings, setBookings] = useState<Booking[]>([])
  const { pushToast } = useToast()

  const loadBookings = async () => {
    try {
      const data = await apiFetch('/api/bookings', { token })
      setBookings(data)
    } catch (err) {
      pushToast('error', (err as Error).message)
    }
  }

  useEffect(() => {
    loadBookings()
  }, [])

  const cancelBooking = async (bookingId: number) => {
    try {
      await apiFetch(`/api/bookings/${bookingId}`, {
        method: 'PATCH',
        token,
        body: JSON.stringify({ status: 'cancelled' }),
      })
      await loadBookings()
      pushToast('success', 'Booking cancelled.')
    } catch (err) {
      pushToast('error', (err as Error).message)
    }
  }

  return (
    <div className="page-body">
      <header className="page-header">
        <div>
          <span className="kicker">My bookings</span>
          <h2>Track your confirmed schedules</h2>
          <p>{user?.username}, review or cancel upcoming reservations.</p>
        </div>
      </header>

      <section className="card">
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
                {booking.status === 'confirmed' && (
                  <button
                    type="button"
                    className="button ghost"
                    onClick={() => cancelBooking(booking.id)}
                  >
                    Cancel
                  </button>
                )}
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
