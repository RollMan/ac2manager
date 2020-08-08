const path = require('path')

module.exports = [{
  mode: "development",
  entry: {
    "index": "./static/assets/js/index.js",
    "login": "./static/assets/js/login.js"
  },
  output: {
    filename: "[name].bundle.js",
    path: path.resolve(__dirname, "static/assets/js")
  },
  devtool: 'inline-source-map'
},
{
  mode: "development",
  entry: {
    "admin": "./admin/assets/js/admin.js",
    "add_race": "./admin/assets/js/add_race.js",
    "edit_race": "./admin/assets/js/edit_race.js"
  },
  output: {
    filename: "[name].bundle.js",
    path: path.resolve(__dirname, "admin/assets/js")
  },
  devtool: 'inline-source-map'
}]
