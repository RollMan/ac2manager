import {table_template} from "/assets/js/event_table_template.js";
(function(){
  const event_table_div = document.querySelector('div#event-tables-div')
  function show_all_events (){
    fetch('/api/races', {
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
        return response.json();
    }).then(response => {
      let body = "";
      for (let e_idx = 0; e_idx < response.length; e_idx++){
        body += table_template(response[e_idx]);
      }
      event_table_div.innerHTML = body
    });
  }

  show_all_events();
})();