styles:
	tailwindcss -i ./static/input.css -o ./static/output.css

run: styles
	go run main.go
