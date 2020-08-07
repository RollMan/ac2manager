import { parseISO, formatRFC3339 } from 'date-fns'

(function(){
  function fetch_upcoming(){
    const upcoming_race_div = document.querySelector('div#upcoming_race')

    fetch("/api/upcoming_race", {
      headers:{
        method: 'GET'
      }
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
          throw new Error("Login Failed.")
        }
        return response.json()
      })
    .then(responses => {
      if(responses.length == 0){
        upcoming_race_div.innerHTML = "<p>No races scheduled.</p>";
      }else{
        let response = responses[0]
        const startdate = parseISO(response.startdate)
        const time_zone = 'Asia/Tokyo'
        const time_jst_string = formatRFC3339(startdate, {'timeZone': time_zone})
        response.startdate = time_jst_string
        delete response.id;
        let body = "<table>"
        for (let field in response) {
          let key = field;
          let value = response[field]
          body += `<tr><td>${key}</td><td>${value}</td></tr>`
        }
        body += '</table>'
        upcoming_race_div.innerHTML = "<p>" + body + "</p>"
      }
    });
  }
  fetch_upcoming();
})();
