package smtp

const emailTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>{{.Subject}}</title>
</head>
<body style="font-family: Arial, sans-serif; margin: 0; padding: 20px; background-color: #f4f4f4;">
    <div style="max-width: 600px; margin: 0 auto; background-color: #ffffff; padding: 20px; border-radius: 8px;">
        <h2 style="color: #333;">{{.AppName}}</h2>
        <h3>{{.Subject}}</h3>
        <p>{{.Body}}</p>
        <p style="color: #777;">Email này được gửi từ {{.AppName}}. Vui lòng không trả lời trực tiếp.</p>
    </div>
</body>
</html>
`