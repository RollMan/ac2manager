function fetchAllEventData(){
}

document.addEventListener('DOMContentLoaded', () => {
  const table_template = document.querySelector('#raceevent-table-template').content;
  const row_template = document.querySelector('#raceeven-raw-template').content;
  const fragment = document.createDocumentFragment();

  const eventData = fetchAllEventData();

  for(const data of eventData){
    const clone = document.importNode(row_template, true);

    const date = clone.querySelector(".date");
    const track = clone.querySelector(".track");
    const R_sessionDurationMinute = clone.querySelector(".race_duration");

    date.textContent = data.startdate;
    track.textContent = data.track;
    R_sessionDurationMinute.textContent = data.R_sessionDurationMinute;

    fragment.appendChild(clone);

  }

  document.querySelector('#raceevent-row-span').appendChild(fragment);
});
