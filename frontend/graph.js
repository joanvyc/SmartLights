// Set the dimensions of the canvas / graph
var margin = {top: 30, right: 20, bottom: 30, left: 50};
var width  = 600 - margin.left - margin.right;
var height = 270 - margin.top - margin.bottom;

// Parse the date / time
var parseDate = d3.time.format("%d-%b-%y-%H:%M.%S").parse;

// Set the ranges
var x = d3.time.scale().range([0, width]);
var y = d3.scale.linear().range([height, 0]);

// Define the axes
var xAxis = d3.svg.axis().scale(x)
    .orient("bottom").ticks(5);

var yAxis = d3.svg.axis().scale(y)
    .orient("left").ticks(5);

var lastDate = null;

// Define the line
var valueline = d3.svg.line()
    .x(function(d) { return x(d.date); })
    .y(function(d) { return y(d.value); })
	.defined(function(d) {
		if ( lastDate === null ) {
			visible = true;
		} else {
			visible = (d.date - lastDate) < (15 * 60 * 1000);
		}
		lastDate = d.date;
		return visible;
		//return true;
	});
    
// Adds the svg canvas
var svg = d3.select("body")
    .append("svg")
        .attr("width", width + margin.left + margin.right)
        .attr("height", height + margin.top + margin.bottom)
    .append("g")
        .attr("transform", 
              "translate(" + margin.left + "," + margin.top + ")");

function flushSvg() {
	d3.select("svg").remove()
	svg = d3.select("body")
    .append("svg")
        .attr("width", width + margin.left + margin.right)
        .attr("height", height + margin.top + margin.bottom)
    .append("g")
        .attr("transform", 
              "translate(" + margin.left + "," + margin.top + ")");
}
// Get the data

var data =
[{date: "10-Apr-12-10:20.12",value: 58.13},
{date: "10-Apr-12-10:24.12",value: 53.98},
{date: "10-Apr-12-10:28.12",value: 67.00},
{date: "10-Apr-12-10:30.12",value: 89.70},
{date: "10-Apr-12-10:31.12",value: 99.00},
{date: "10-Apr-12-10:34.12",value: 130.28},
{date: "10-Apr-12-10:39.12",value: 166.70},
{date: "10-Apr-12-10:47.12",value: 234.98},
{date: "10-Apr-12-10:48.12",value: 345.44},
{date: "10-Apr-12-11:30.12",value: 443.34},
{date: "10-Apr-12-12:25.12",value: 543.70},
{date: "10-Apr-12-12:25.12",value: 580.13},
{date: "10-Apr-12-12:26.12",value: 605.23},
{date: "10-Apr-12-12:27.12",value: 622.77},
{date: "10-Apr-12-12:29.12",value: 626.20},
{date: "10-Apr-12-12:32.12",value: 628.44},
{date: "10-Apr-12-12:35.12",value: 636.23},
{date: "10-Apr-12-12:36.12",value: 633.68},
{date: "10-Apr-12-12:38.12",value: 624.31},
{date: "10-Apr-12-12:44.12",value: 629.32},
{date: "10-Apr-12-12:46.12",value: 618.63}]
data.forEach(function(d) { d.date = parseDate(d.date);})

function draw(data) {
	x.domain(d3.extent(data, function(d) { return d.date; }));
	y.domain([0, d3.max(data, function(d) { return d.value; })]);

	// Add the valueline path.
	svg.append("path")
		.attr("class", "line")
		.attr("d", valueline(data));

	// Add the X Axis
	svg.append("g")
		.attr("class", "x axis")
		.attr("transform", "translate(0," + height + ")")
		.call(xAxis);

	// Add the Y Axis
	svg.append("g")
		.attr("class", "y axis")
		.call(yAxis);
}

var socket = new WebSocket("ws://192.168.1.118:3000/ws");

socket.onopen = function(event) {
	console.log("Socket opened.");
}

socket.onerror = function(error) {
	console.log("Something went wrong! Log: " + error.message);
}

socket.onmessage = function(event) {
	socketData = parse (event.data)
}


draw(data);
setTimeout(function() {
	data.push({date: parseDate("11-Apr-12-13:46.12"),value: 700.32})
	data.push({date: parseDate("11-Apr-12-14:00.12"),value: 707.32})
	data.push({date: parseDate("11-Apr-12-14:10.12"),value: 719.32})
	data.push({date: parseDate("11-Apr-12-14:20.12"),value: 722.32})
	data.push({date: parseDate("11-Apr-12-14:30.12"),value: 745.32})
	flushSvg();
	draw(data);
}, 2000);
