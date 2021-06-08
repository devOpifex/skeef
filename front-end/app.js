var createGraph = require('ngraph.graph');
var pixel = require('ngraph.pixel');

let trendPlot,
    graphContainer,
    clickedNode;

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
  graphContainer = document.getElementById("graph");

  if(graphContainer != null){
    var renderer = pixel(g, {
      is3d: true,
      container: graphContainer,
      clearAlpha: 0.0,
      autoFit: true,
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
            size: sizeNode(n.data.count + 1)
          }
        } else if(n.data.type == 'hashtag'){
          return {
            color: 0x9d4edd,
            size: sizeNode(n.data.count + 1)
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
          fromColor: 0xbfbfbf,
          toColor: 0xbfbfbf
        }
      }
    });

    renderer.on('nodeclick', function(node){
      clickedNode = node.id;
      socket.send(node.id);
    })

    resizeGraph();
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

  let firstrun = true,
      nnodes = 0,
      nedges = 0;
  let nedgesEl = document.getElementById("nedges");
  let nnodesEl = document.getElementById("nnodes");
  trendPlot = initPlot();
  let tweetsWrapper = document.getElementById("tweets");

  window.socket.onmessage = (data) => {
    if(nedgesEl == null)
      return;

    let parsed = JSON.parse(data.data);

    if(parsed != null && parsed.hasOwnProperty("embeds")){
      // empty the wrapper
      tweetsWrapper.innerHTML = "";

      if(parsed.embeds == null){
        return ;
      }

      let embeds = parsed.embeds.filter(uniques);

      for(let i = 0; i < embeds.length; i++){

        // create div
        let el = document.createElement("DIV");
        tweetsWrapper.appendChild(el);

        // embed tweet 
        twttr.widgets.createTweet(
          embeds[i], el, {
            theme: 'dark',
            conversation: 'none',
            cards: 'hidden'
          }
        )

        setTimeout(function(){
          document.getElementById('hideTweets').style.display = 'block';
        }, 500);
      }
      return;
    }

    updatePlot(parsed.trend);
    
    // Initial load, all at once
    if(firstrun){
      g.beginUpdate();
      for(let i = 0; i < parsed.graph.nodes.length; i++) {
        nnodes++
        g.addNode(
          parsed.graph.nodes[i].name,
          {type: parsed.graph.nodes[i].type, count: parsed.graph.nodes[i].count}
        );
      }
      for(let i = 0; i < parsed.graph.edges.length; i++) {
        nedges++
        g.addLink(parsed.graph.edges[i].source, parsed.graph.edges[i].target);
      }
      g.endUpdate();
      firstrun = false
      return ;
    }

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
          nnodes++
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
          nedges++
          g.addLink(parsed.graph.edges[i].source, parsed.graph.edges[i].target);
        }
        
      }
    }
    g.endUpdate();

    let lnk;
    if(parsed.removeEdges != null){
      g.beginUpdate();
      for( let i = 0; i < parsed.removeEdges.length; i++){
        nedges--
        lnk = g.getLink(
          parsed.removeEdges[i].source,
          parsed.removeEdges[i].target
        )
        g.removeLink(lnk);
      }
      g.endUpdate();
    }

    if(parsed.removeNodes != null){
      g.beginUpdate();
      for( let i = 0; i < parsed.removeNodes.length; i++){
        nnodes--
        g.removeNode(parsed.removeNodes[i].name);
      }
      g.endUpdate();
    }

    if(parsed.graph.nodes != null){
      nnodesEl.innerText = nnodes;
    }

    if(parsed.graph.edges != null){
      nedgesEl.innerText = nedges;
    }

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
    width: window.innerWidth - 40,
    height: 120,
    title: "",
    axes: [
      {
        stroke: "#c7d0d9",
      //	font: `12px 'Roboto'`,
      //	labelFont: `12px 'Roboto'`,
        grid: {
          show: false
        },
        ticks: {
          width: 1 / devicePixelRatio,
          stroke: "#2c3235",
        }
      },
      {
        stroke: "#c7d0d9",
      //	font: `12px 'Roboto'`,
      //	labelFont: `12px 'Roboto'`,
        grid: {
          width: 1 / devicePixelRatio,
          stroke: "#2c3235",
        },
        ticks: {
          width: 1 / devicePixelRatio,
          stroke: "#2c3235",
        }
      },
    ],
    scales: {
      x: {
        time: true,
      },
    },
    series: [
      {
        label: 'Date time'
      },
      {
        label: 'Tweets',
        stroke: "#ff9e00",
        fill: "rgba(255,158,0,0.1)",
        width: 1/devicePixelRatio,
      }
    ],
  };

  let u = new uPlot(opts, data, document.getElementById("trend"));

  return u;
}

function updatePlot(data){
  if(data == undefined)
    return ;

  if(data == null){
    return ;
  }

  if(Object.keys(data).length < 2){
    return ;
  }

  let result = [
    Object.keys(data),
    Object.values(data)
  ]

  trendPlot.setData(result);
}

function getSize() {
  return {
    width: window.innerWidth - 40,
    height: 120,
  }
}

window.addEventListener("resize", e => {
  trendPlot.setSize(getSize());
  resizeGraph();
});

function resizeGraph(){
  if(graphContainer == null)
    return ;

  graphContainer.childNodes[0].style.width = graphContainer.offsetWidth + 'px';
}

function sizeNode(n){
  return Math.round(Math.log10(n + 1) * 25);
}

function uniques(value, index, self) {
  return self.indexOf(value) === index;
}