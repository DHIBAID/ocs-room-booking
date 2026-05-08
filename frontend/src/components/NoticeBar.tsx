type NoticeBarProps = {
  tone: 'neutral' | 'error'
  message: string
}

export const NoticeBar = ({ tone, message }: NoticeBarProps) => {
  if (!message) {
    return null
  }

  return (
    <div className={`notice ${tone === 'error' ? 'notice-error' : ''}`}>
      {message}
    </div>
  )
}
