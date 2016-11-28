"use strict";

/**
 * Fetches new data and updates the page.
 */
function update() {
  var req = new XMLHttpRequest();
  req.onreadystatechange = function() {
    // Update the page
    if (this.status == 200 && this.responseText != "") {
      updateGraphs(JSON.parse(this.responseText));
    }
  };
  req.open("GET", "/data", true);
  req.send();
}

function updateGraphs(data) {
  console.log("updating");

  var loading = document.getElementById("loading");
  var content = document.getElementById("content");

  window.accelGraph.config.data.datasets[0].data = data.accel;
  window.accelGraph.config.data.labels = data.accel;
  window.accelGraph.update();
  window.gyroGraph.config.data.datasets[0].data = data.gyro;
  window.accelGraph.config.data.labels = data.gyro;
  window.gyroGraph.update();

  loading.style.display='none';
}

document.addEventListener("DOMContentLoaded", function() {
  // Initialize the graphs
  var accelGraphCtx = document.getElementById("accelGraph").getContext("2d");
  var gyroGraphCtx = document.getElementById("gyroGraph").getContext("2d");
  window.accelGraph = new Chart(accelGraphCtx, accelGraphConfig);
  window.accelGraph.update();
  window.gyroGraph = new Chart(gyroGraphCtx, gyroGraphConfig);
  window.gyroGraph.update();

  // Update the page every second
  setInterval(update, 1000);

}, false);

var userDataColor = 'rgba(174, 54, 90, 0.5)';
var correctDataColor = 'rgba(54, 174, 90, 0.5)';

var accelGraphConfig = {
  type: 'line',
  data: {
    labels: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
    datasets: [{
      fill: true,
      fillColor: userDataColor,
      pointColor: userDataColor,
      backgroundColor: userDataColor,
      data: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
    }]
  },
  options: {
    responsive: true,
    maintainAspectRatio: false,
    legend: {
      display: false,
    },
    title: {
      display: false,
    },
    scales: {
      xAxes: [{
        display: false,
        scaleLabel: {
          display: false,
        }
      }],
      yAxes: [{
        display: false,
        scaleLabel: {
          display: false,
        }
      }]
    }
  }
};

var gyroGraphConfig = JSON.parse(JSON.stringify(accelGraphConfig));

