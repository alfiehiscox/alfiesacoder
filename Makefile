templ: 
	templ generate

styles: templ
	tailwindcss -i ./static/input.css -o ./static/output.css

run: styles
	go run *.go
