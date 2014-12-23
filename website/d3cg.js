/* Copyright (c) 2014 Marek Bardoński */
(function () {

    var svg = d3.select("main")
            .append("svg"),
        margin = {top: 10,
                  right: 10,
                  bottom: 10,
                  left: 10},
        parse = function (N) {
            return Number(N.replace("px", ""));
        },
        currentKeysPressed = [];
        playerDX = 0;
        playerDY = 0;
        targetDX = 0;
        targetDY = 0;
    
    d3.select("body").on("contextmenu", function(data, index) {
		var coordinates = [0, 0];
		coordinates = d3.mouse(this);
		var x = coordinates[0];
		var y = coordinates[1];
		
		targetDX = x - Screen().width / 2;
		targetDY = y - Screen().height / 2;
		
		//stop showing browser menu
		d3.event.preventDefault();
	});
        
    // Add support for movement by keys
    // When a key is pressed, add it to the current keys array for further tracking
    d3.select("body").on("keydown", function() {
        if (currentKeysPressed.indexOf(d3.event.keyCode) != -1) { return }
        currentKeysPressed.push(d3.event.keyCode);
    });

    // When the key is relased, remove it from the array.
    d3.select("body").on("keyup", function() {
        currentKeysPressed.splice(currentKeysPressed.indexOf(d3.event.keyCode), 1);
    });

    // always returns current SVG dimensions
    var Screen = function () {
            return {
                width: parse(svg.style("width")),
                height: parse(svg.style("height"))
            };
    },
    
       Ball = function () {
            var R = 5,
                ball = svg.append('circle')
                    .classed("ball", true)
                    .attr({r: R,
                           cx: Screen().width/2 + playerDX,
                           cy: Screen().height/2 + playerDY}),
                scale = d3.scale.linear().domain([0, 1]).range([-1, 1]),
                vector = {x: scale(Math.random()),
                          y: scale(Math.random())},
                speed = 7;
			
			return function f(left, right, delta_t) {
                var screen = Screen(),
                    // this should pretend we have 100 fps
                    fps = delta_t > 0 ? (delta_t/1000)/100 : 1; 

                ball.attr({
                    cx: parse(ball.attr("cx"))+vector.x*speed*fps,
                    cy: parse(ball.attr("cy"))+vector.y*speed*fps
                });

                return false;
            };
		};
	
	ball = Ball();
   
    // detect window resize events (also captures orientation changes)
    d3.select(window).on('resize', function () {
        var screen = Screen();
        d3.select(".ball").remove();
        var ball = Ball();
    });
    
    var i = 0;
    
    function particle() {
		svg.append("circle")
			  .attr("cx", d3.select(".ball").attr("cx"))
			  .attr("cy",d3.select(".ball").attr("cy"))
			  .attr("r", 1e-6)
			  .style("stroke", d3.hsl((i = (i + 1) % 360), 1, .5))
			  .style("stroke-opacity", 1)
			.transition()
			  .duration(2000)
			  .ease(Math.sqrt)
			  .attr("r", 100)
			  .style("stroke-opacity", 1e-6)
			  .remove();
	 }
    
    function processKeys() {
        for (var i = 0; i < currentKeysPressed.length; i++) {
            var currentKeyPressed = currentKeysPressed[i];
            
            var speed = 5;
            
            // Q
            if (currentKeyPressed == 81) {
				particle();
			}
			
			// W 
			if (currentKeyPressed == 87) {
			}
			
			// E
			if (currentKeyPressed == 69) {
			}
			
			// R 
			if (currentKeyPressed == 82) {
			}
            
            if (currentKeyPressed == 38) {
				playerDY -= speed;
			}
			
            if (currentKeyPressed == 40) {
				playerDY += speed;
			}

            if (currentKeyPressed == 37) {
				playerDX -= speed;
			}
			
            if (currentKeyPressed == 39) {
				playerDX += speed;
			}			
        }
    }
    
    function movePlayer(){
		if (playerDX != targetDX) {
			playerDX += Math.max(-3, Math.min(3, targetDX - playerDX));
		}
		
		if (playerDY != targetDY) {
			playerDY += Math.max(-3, Math.min(3, targetDY - playerDY));
		}

		var player = d3.select('.ball');
		player
			.transition()
			.attr('cx', Screen().width/2 + playerDX)
			.attr('cy', Screen().height/2 + playerDY)
			.duration(30);
	}

    function run() {
        var last_time = Date.now();
        d3.timer(function () {
			
			var now = Date.now();
			
			/// 30FPS
			if (now-last_time > 30) {
				processKeys();
				movePlayer();
				last_time = now;
			}
			
			// Infinite loop
			return false;
        }, 50);
    };

    run();
    
})();
