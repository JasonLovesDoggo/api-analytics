import { join } from "path";
import express from "express";
import app from "./public/App.js";
import { dirname } from "path";
import { fileURLToPath } from "url";

const server = express();

const __dirname = dirname(fileURLToPath(import.meta.url));

server.use(express.static(join(__dirname, "public")));

server.get("*", function (req, res) {
  const { html } = app.render({ url: req.url });

  res.write(`
    <!DOCTYPE html>
    <head>
      <link rel='stylesheet' href='/global.css'>
      <link rel='stylesheet' href='/bundle.css'>
      <link rel="icon" type="image/x-icon" href="/img/favicon.ico">
      <title>API Analytics</title>
      <script src="https://cdn.plot.ly/plotly-latest.min.js" type="text/javascript"></script>
      <meta name="viewport" content="width=device-width, initial-scale=1.0">
    </head>

    <body>
      <div id="app">${html}</div>
      <script src="/bundle.js"></script>
    </body>
  `);

  res.end();
});

const port = 3000;
server.listen(port, () => console.log(`Listening on http://localhost:${port}`));

