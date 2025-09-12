package responses

import (
	"fmt"
	"html"
	"net/http"
	"os"
)

// ErrorHTML renders a centered card with an error message, using inline styles only.
func ErrorHTML(w http.ResponseWriter, msg string) {
	title := "Error"
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)

	// Base URL for your logo
	complyageLogo := os.Getenv("COMPLYAGE_CLIENT_URL")

	safeTitle := html.EscapeString(title)
	safeMsg := html.EscapeString(msg)

	fmt.Fprintf(w, `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>%s</title>
</head>
<body style="margin:0;height:100vh;display:flex;align-items:center;justify-content:center;background:#f0f2f5;font-family:Arial,sans-serif;">
  <div style="
      background:#fff;
      padding:24px;
      border-radius:8px;
      box-shadow:0 4px 16px rgba(0,0,0,0.1);
      text-align:center;
      max-width:360px;
      width:90%%;
  ">
    <img
      src="%s/static/media/complyage.webp"
      alt="Logo"
    />
    <h1 style="margin:0 0 12px;font-size:22px;color:#e55353;">%s</h1>
    <p style="font-weight:bold;margin:0;font-size:16px;color:#555;line-height:1.4;">%s</p>
  </div>
</body>
</html>`,
		safeTitle,     // <title>
		complyageLogo, // logo base URL
		safeTitle,     // H1 text
		safeMsg,       // paragraph text
	)
}
