import { parseISO, format } from 'date-fns'

(function(){
  const event_div = document.querySelector("div#event_edit_table");
  const submit_button = document.querySelector("button#submit_race");
  const edit_form = document.querySelector("form#edit_race")
  const edit_result_div = document.querySelector("div#event_edit_result")
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

  submit_race.addEventListener("click", function(){
    const form_html_element = document.querySelector("form#edit_race")
    const form_data = new FormData(form_html_element);

    const date = form_data.get("startdatedate")
    const time = form_data.get("startdatetime")
    form_data.delete("startdatedate")
    form_data.delete("startdatetime")

    const datetime_rfc3339 = date + "T" + time + ":00+09:00";
    form_data.append("startdate", datetime_rfc3339);

    form_data.append("id", id);

    fetch('/api/remove_race', {
      method: 'POST',
      headers: {"Content-Type": "application/json"},
      body: `{"id": ${id}}`
    })
    .then(response => {
        if (!response.ok) {
          let body = "";
          if (400 <= response.status && response.status < 500){
            body = "<p>Client Error</p>"
          }else if(500 <= response.status && response.status < 600){
            body = "<p>Internal Server Error.</p>"
          }else{
            body = "<p>Unknown Error: " + response.status + ".</p>"
          }
          response.text().then(reason => {
            edit_result_div.innerHTML = body + "<br>" + reason;
          });
          throw new Error("Remove race followed by add race failed")
        }
      return response.json()
    })
    .then(json => {
      const removed_count = json.count;
      if (removed_count != 1) {
        let body = "";
        if (removed_count == 0) {
          body = `<p>Internal Server Error: id ${id} does not exist.</p>`
        }else if (removed_count > 1) {
          body = `<p>Internal Server Error: id ${id} corresponds to multiple entries and they were removed. Contact the administrator immediately.</p>`;
        }
        edit_result_div.innerHTML = body;
        throw new Error("Remove Failed.")
      }
      fetch('/api/add_race', {
        method: 'POST',
        body: form_data,
      })
        .then(response => {
          if (!response.ok) {
            let body = "";
            if (400 <= response.status && response.status < 500){
              body = "<p>Client Error</p>"
            }else if(500 <= response.status && response.status < 600){
              body = "<p>Internal Server Error.</p>"
            }else{
              body = "<p>Unknown Error: " + response.status + ".</p>"
            }
            response.text().then(reason => {
              edit_result_div.innerHTML = body + "<br>" + reason;
            });
            throw new Error("Add race failed")
          }
          return response.json()
        })
        .then(response => {
          edit_result_div.innerHTML = "ok";
        })
    });
    });

})();
