import { FormEvent, useEffect, useMemo, useState } from 'react'
import Select from 'react-select'
import { apiFetch } from '../api/client'
import type { Room } from '../api/types'
import { useAuth } from '../hooks/useAuth'
import { useToast } from '../components/ToastProvider'

const purposes = ['OA', 'Interview', 'PPT']

export const DashboardPage = () => {
  const { token } = useAuth()
  const today = new Date().toISOString().split('T')[0]
  const [search, setSearch] = useState({
    date: '',
    start_time: '',
    end_time: '',
    purpose: 'OA',
    minCapacity: '',
    block: '',
  })
  const [rooms, setRooms] = useState<Room[]>([])
  const [allRooms, setAllRooms] = useState<Room[]>([])
  const [bookingForm, setBookingForm] = useState({
    room_id: null as number | null,
    participants: '',
    purpose: 'OA',
  })
  const { pushToast } = useToast()

  const roomOptions = useMemo(() => {
    const source = allRooms.length > 0 ? allRooms : rooms
    return source.map((room) => ({
      value: room.id,
      label: room.name,
      meta: `${room.block} · ${room.capacity} seats`,
    }))
  }, [allRooms, rooms])

  useEffect(() => {
    apiFetch('/api/rooms', { token })
      .then((data) => setAllRooms(data))
      .catch((err) => pushToast('error', (err as Error).message))
  }, [token, pushToast])

  const handleSearch = async (event: FormEvent) => {
    event.preventDefault()
    if (search.date && search.date < today) {
      pushToast('error', 'Past dates are not allowed.')
      return
    }
    const params = new URLSearchParams()
    if (search.block) params.set('block', search.block)
    if (search.purpose) params.set('purpose', search.purpose)
    if (search.minCapacity) params.set('minCapacity', search.minCapacity)
    if (search.date) params.set('date', search.date)
    if (search.start_time) params.set('start_time', search.start_time)
    if (search.end_time) params.set('end_time', search.end_time)

    try {
      const data = await apiFetch(`/api/rooms?${params.toString()}`, { token })
      setRooms(data)
      pushToast('success', `Loaded ${data.length} rooms.`)
    } catch (err) {
      pushToast('error', (err as Error).message)
    }
  }

  const handleCreateBooking = async (event: FormEvent) => {
    event.preventDefault()
    if (!search.date || !search.start_time || !search.end_time) {
      pushToast('error', 'Choose a date and time window before booking.')
      return
    }
    if (search.date < today) {
      pushToast('error', 'Bookings cannot be made for past dates.')
      return
    }

    const selectedRoom = (allRooms.length > 0 ? allRooms : rooms).find(
      (room) => room.id === bookingForm.room_id,
    )
    if (!selectedRoom || bookingForm.room_id === null) {
      pushToast('error', 'Select a room from the availability list.')
      return
    }

    try {
      await apiFetch('/api/bookings', {
        method: 'POST',
        token,
        body: JSON.stringify({
          room_id: selectedRoom.id,
          participants: Number(bookingForm.participants),
          purpose: bookingForm.purpose,
          date: search.date,
          start_time: search.start_time,
          end_time: search.end_time,
        }),
      })
      pushToast('success', 'Booking confirmed. Check the bookings page.')
    } catch (err) {
      pushToast('error', (err as Error).message)
    }
  }

  return (
    <div className="page-body">
      <header className="page-header">
        <div>
          <span className="kicker">Availability</span>
          <h2>Find rooms that fit your slot</h2>
          <p>Search by block, capacity, and purpose to reserve quickly.</p>
        </div>
      </header>

      <div className="page-grid">
        <section className="card card-merge">
          <div className="card-merge__header">
            <h3>Search & confirm booking</h3>
            <p>Use a single flow to find rooms and lock in the reservation.</p>
          </div>
          <div className="card-merge__grid">
            <form className="form" onSubmit={handleSearch}>
              <h4>Search inventory</h4>
              <div className="field-grid">
                <label className="field">
                  <span>Date</span>
                  <input
                    className="input"
                    type="date"
                    min={today}
                    value={search.date}
                    onChange={(event) =>
                      setSearch({ ...search, date: event.target.value })
                    }
                  />
                </label>
                <label className="field">
                  <span>Start</span>
                  <input
                    className="input"
                    type="time"
                    value={search.start_time}
                    onChange={(event) =>
                      setSearch({ ...search, start_time: event.target.value })
                    }
                  />
                </label>
                <label className="field">
                  <span>End</span>
                  <input
                    className="input"
                    type="time"
                    value={search.end_time}
                    onChange={(event) =>
                      setSearch({ ...search, end_time: event.target.value })
                    }
                  />
                </label>
              </div>
              <div className="field-grid">
                <label className="field">
                  <span>Purpose</span>
                  <select
                    className="input"
                    value={search.purpose}
                    onChange={(event) =>
                      setSearch({ ...search, purpose: event.target.value })
                    }
                  >
                    {purposes.map((purpose) => (
                      <option key={purpose} value={purpose}>
                        {purpose}
                      </option>
                    ))}
                  </select>
                </label>
                <label className="field">
                  <span>Minimum capacity</span>
                  <input
                    className="input"
                    type="number"
                    min="1"
                    value={search.minCapacity}
                    onChange={(event) =>
                      setSearch({
                        ...search,
                        minCapacity: event.target.value,
                      })
                    }
                  />
                </label>
                <label className="field">
                  <span>Block</span>
                  <input
                    className="input"
                    value={search.block}
                    onChange={(event) =>
                      setSearch({ ...search, block: event.target.value })
                    }
                  />
                </label>
              </div>
              <button className="button" type="submit">
                Check availability
              </button>
            </form>

            <form className="form" onSubmit={handleCreateBooking}>
              <h4>Confirm booking</h4>
              <div className="field-grid">
                <label className="field">
                  <span>Date</span>
                  <input
                    className="input"
                    type="date"
                    min={today}
                    value={search.date}
                    onChange={(event) =>
                      setSearch({ ...search, date: event.target.value })
                    }
                  />
                </label>
                <label className="field">
                  <span>Start</span>
                  <input
                    className="input"
                    type="time"
                    value={search.start_time}
                    onChange={(event) =>
                      setSearch({ ...search, start_time: event.target.value })
                    }
                  />
                </label>
                <label className="field">
                  <span>End</span>
                  <input
                    className="input"
                    type="time"
                    value={search.end_time}
                    onChange={(event) =>
                      setSearch({ ...search, end_time: event.target.value })
                    }
                  />
                </label>
              </div>
              <label className="field">
                <span>Selected room name</span>
                <Select
                  classNamePrefix="room-select"
                  className="room-select"
                  placeholder="Pick a room from the list"
                  options={roomOptions}
                  value={roomOptions.find(
                    (option) => option.value === bookingForm.room_id,
                  )}
                  onChange={(option) =>
                    setBookingForm({
                      ...bookingForm,
                      room_id: option ? option.value : null,
                    })
                  }
                  noOptionsMessage={() =>
                    rooms.length === 0
                      ? 'Run a search to load rooms.'
                      : 'No matching rooms.'
                  }
                  formatOptionLabel={(option) => (
                    <div className="room-option-inline">
                      <span>{option.label}</span>
                      <span className="room-option-inline__meta">
                        {option.meta}
                      </span>
                    </div>
                  )}
                  isClearable
                />
              </label>
              <div className="field-grid">
                <label className="field">
                  <span>Participants</span>
                  <input
                    className="input"
                    type="number"
                    min="1"
                    value={bookingForm.participants}
                    onChange={(event) =>
                      setBookingForm({
                        ...bookingForm,
                        participants: event.target.value,
                      })
                    }
                    required
                  />
                </label>
                <label className="field">
                  <span>Purpose</span>
                  <select
                    className="input"
                    value={bookingForm.purpose}
                    onChange={(event) =>
                      setBookingForm({
                        ...bookingForm,
                        purpose: event.target.value,
                      })
                    }
                  >
                    {purposes.map((purpose) => (
                      <option key={purpose} value={purpose}>
                        {purpose}
                      </option>
                    ))}
                  </select>
                </label>
              </div>
              <button className="button" type="submit">
                Submit booking
              </button>
            </form>
          </div>
        </section>

        <section className="card rooms-panel">
          <h3>Available rooms</h3>
          <div className="room-grid">
            {rooms.map((room) => (
              <button
                type="button"
                key={room.id}
                className={`room-card ${
                  bookingForm.room_name === room.name ? 'active' : ''
                }`}
                onClick={() =>
                  setBookingForm({
                    ...bookingForm,
                    room_name: room.name,
                  })
                }
              >
                <div className="room-title">
                  <span className="room-badge">{room.block}</span>
                  <span>{room.name}</span>
                </div>
                <div className="room-meta">
                  <span>{room.capacity} seats capacity</span>
                  <span>{room.status}</span>
                  {room.allowed_purposes && (
                    <span>Allowed: {room.allowed_purposes}</span>
                  )}
                  {room.notes && <span>Notes: {room.notes}</span>}
                </div>
              </button>
            ))}
            {rooms.length === 0 && (
              <p className="empty">Run a search to load available rooms.</p>
            )}
          </div>
        </section>

      </div>
    </div>
  )
}
