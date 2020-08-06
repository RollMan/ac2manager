(function(){
  const event_table_for_edit_div = document.querySelector('div#event_table_for_edit')
  const result_div = document.querySelector('div#event_add_result')
  const submit_button = document.querySelector('button#submit_race')


  fetch('/api/schema', {
    method: 'GET',
    headers: {'Content-Type': 'application/json'}
  })
    .then(response => {
      if (!response.ok) {
        let body = "";
        if (400 <= response.status && response.status < 500){
          body = "<p>Client Error.</p>"
        }else if(500 <= response.status && response.status < 600){
          body = "<p>Internal Server Error.</p>"
        }else{
          body = "<p>Unknown Error: " + response.status + ".</p>"
        }
        upcoming_race_div.innerHTML = body;
        throw new Error("Server Unavailable.")
      }
      return response.json()
    })
    .then(response => {
      let schema = response;
      let body = "<table>";
      for (let field in schema) {
        let key = field;
        let type = schema[field];
        let row = `<tr><td>${key}</td>`;
        if (key == "id") continue;
        if (type == "int" || type == "uint"){
          row += `<td><input type="number" id="${key}" name="${key}"></td>`;
        }else if (type == "string"){
          row += `<td><input type="text" id="${key}" name="${key}"></td>`;
        }else if (type == "bool"){
          row += `<td><input type="checkbox" id="${key}" name="${key}"></td>`;
        }else if (type == "Time"){
          row += `<td><input type="date" id="${key}date" name="${key}date"><br><input type="time" id="${key}time" name="${key}time"></td>`;
        }
        body += row + "</tr>";
      }
      body += "</table>"
      event_table_for_edit_div.innerHTML = body;
    });
})();
