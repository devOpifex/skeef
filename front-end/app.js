var createGraph = require('ngraph.graph');
var pixel = require('ngraph.pixel');

document.addEventListener("DOMContentLoaded",function(){

  // warning on admin inputs
  let locationsInput = document.getElementById("locationsStream");
  
  if(locationsInput != null){
    let locationsWarning = document.getElementById("locationsWarning");
    locationsInput.addEventListener("input", (event) => {
      if(event.data == null){
        locationsWarning.innerHTML = "";
        return ;
      }
      locationsWarning.innerHTML = "Few tweets are geo-tagged; this dramatically reduces the number of tweets streamed."
    });
  }

  var g = createGraph();
  g.addNode(1, {type: 'hidden'});
  g.addNode(2, {type: 'hidden'});

  if(document.getElementById("graph") != null){
    var renderer = pixel(g, {
      is3d: true,
      container: document.getElementById("graph"),
      clearAlpha: 0.0,
      physics: {
        springLength : 50,
        springCoeff : 0.0002,
        gravity: -2,
        theta : 0.4,
        dragCoeff : 0.02
      },
      node: function (n) {
        if(n.data.type == 'user'){
          return {
            color: 0xff9e00,
            size: n.data.count * 10
          }
        } else if(n.data.type == 'hashtag'){
          return {
            color: 0x9d4edd,
            size: n.data.count * 10
          }
        } else if (n.data.type == 'hidden') {
          return {
            color: 0x262b36,
            size: 0
          }
        } 
      },
      link: function(l){
        return {
          fromColor: 0xCCCCCC,
          toColor: 0xCCCCCC
        }
      }
    });
  }

  // websocket
  window.socket.onopen = () => {
    socket.send("Client Connected")
  }
  
  window.socket.onclose = () => {
    console.log("Websocket closed")
  }
  
  window.socket.onerror = (error) => {
    console.log(error);
  }

  let firstrun = true;
  let ntweets = document.getElementById("ntweets");
  let trendData = [];
  let trendPlot = initPlot();

  window.socket.onmessage = (data) => {
    let parsed = JSON.parse(data.data);

    trendData.push(prepTrendData(parsed.trend));
    trendPlot.setData(trendData);
    console.log(trendData);
    
    if(parsed.tweetsCount == 0){
      return ;
    }

    // Initial load, all at once
    if(firstrun){
      g.beginUpdate();
      for(let i = 0; i < parsed.graph.nodes.length; i++) {
        g.addNode(
          parsed.graph.nodes[i].name,
          {type: parsed.graph.nodes[i].type, count: parsed.graph.nodes[i].count}
        );
      }
      for(let i = 0; i < parsed.graph.edges.length; i++) {
        g.addLink(parsed.graph.edges[i].source, parsed.graph.edges[i].target);
      }
      g.endUpdate();
      firstrun = false
      return ;
    }

    ntweets.innerText = parsed.tweetsCount;

    g.beginUpdate();
    if(parsed.graph.nodes){
      for(let i = 0; i < parsed.graph.nodes.length; i++)  {

        if(parsed.graph.nodes[i].action == "update"){
          let n = renderer.getNode(parsed.graph.nodes[i].name);

          // if undefined then we have an issue and we must add 
          // the node instead
          if(n == undefined) {
            g.addNode(
              parsed.graph.nodes[i].name,
              {
                type: parsed.graph.nodes[i].type, 
                count: parsed.graph.nodes[i].count
              }
            );
            continue
          }
          n.size = parsed.graph.nodes[i].count;
        }

        if(parsed.graph.nodes[i].action == "add"){
          g.addNode(
            parsed.graph.nodes[i].name,
            {type: parsed.graph.nodes[i].type, count: parsed.graph.nodes[i].count}
          );
        }

      }
    }

    if(parsed.graph.edges){
      for(let i = 0; i < parsed.graph.edges.length; i++){
        if(parsed.graph.edges[i].action == "update"){
          continue; 
        }

        if(parsed.graph.edges[i].action == "add"){
          g.addLink(parsed.graph.edges[i].source, parsed.graph.edges[i].target);
        }
        
      }
    }
    g.endUpdate();

  }

});

function initPlot(){
  let xs = [];
  let vals = [];

  let data = [
    xs,
    vals,
  ];

  const opts = {
    zDate: ts => uPlot.tzDate(new Date(ts * 1e3), tz),
    width: 1920,
    height: 600,
    title: "Area Fill",
    scales: {
      x: {
        time: true,
      },
    },
    series: [
      {},
      {
        stroke: "#ff9e00",
        fill: "rgba(255,158,0,0.1)",
      }
    ],
  };

  let u = new uPlot(opts, data, document.getElementById("trend"));

  return u;
}

function prepTrendData(data){
  if(data == undefined)
    return [];

  if(data == null){
    return [];
  }

  // dates
  let dates = Object.keys(data);
  dates = dates.map((v) => {Date.parse(v) / 1000});

  // trend
  let trend = Object.values(data);

  return [dates, trend];
}