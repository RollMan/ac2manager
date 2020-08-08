import { parseISO, format } from 'date-fns'

(function(){
  const event_div = document.querySelector("div#event_edit_table");
  const query_string = window.location.search;

  const query = [...new URLSearchParams(query_string).entries()].reduce((obj, e) => ({...obj, [e[0]]: e[1]}), {});

  const id = query.id;

  const params = {id: id};
  const qs = new URLSearchParams(params)
  fetch(`/api/race_by_id?${qs}`, {
    method: "GET",
  })
  .then(response =>{
    if (!response.ok) {
      let body = "";
      if (400 <= response.status && response.status < 500){
        body = "<p>Client Error.</p>"
      }else if(500 <= response.status && response.status < 600){
        body = "<p>Internal Server Error.</p>"
      }else{
        body = "<p>Unknown Error: " + response.status + ".</p>"
      }
      event_div.innerHTML = body;
      throw new Error("Error on race_by_id")
    }
    return response.json()
  })
  .then(ev=> {
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
        event_div.innerHTML = body;
        throw new Error("Server Unavailable.")
      }
      return response.json()
    })
    .then(schema => {
      let body = "<table>";
      for (let field in schema){
        let key = field;
        let type = schema[field];
        let value = ev[key];
        let row = "<tr><td>" + key + "</td>";
        if (type == "int" || type == "uint"){
          row += `<td><input type="number" id="${key}" name="${key}" value="${value}"></td>`;
        }else if (type == "string"){
          row += `<td><input type="text" id="${key}" name="${key}" value="${value}"></td>`;
        }else if (type == "bool"){
          let checked = value == true ? "checked" : "";
          row += `<td><input type="checkbox" value="true" id="${key}" name="${key}" ${checked}></td>`;
        }else if (type == "Time"){
          const time_zone = 'Asia/Tokyo'
          let datetime = parseISO(value);
          let date = format(datetime, "yyyy-MM-dd", { timeZone: time_zone});
          let time = format(datetime, "HH:mm", { timeZone: time_zone});
          row += `<td><input type="date" id="${key}date" name="${key}date" value="${date}"><br><input type="time" id="${key}time" name="${key}time" value="${time}"> +09:00</td>`;
        }
        row += "</tr>"
        body += row;
      }
      body += "</table>"
      event_div.innerHTML = body;
    })
  })
})();
