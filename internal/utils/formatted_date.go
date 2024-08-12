package utils

import (
  "time"
)

func FormattedDate() string {
  currentTime := time.Now()
  formattedDateTime := currentTime.Format("02/01/2006 15:04:05")
  return formattedDateTime
}
