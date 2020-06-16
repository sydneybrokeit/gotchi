var current_food = 2;
var max_food = 4;
var current_love = 2;
var max_love = 4;

// create a string used for doing bars
function createBarString(current, max) {
    if (current > max) {
        current = max;
    }
    var output = "";
    output = output.padStart(current, "1").padEnd(max, "0");
    return output;
};


Stage({
    name: 'foodbar',
    image: {
        src : './assets/foodbar.png',
    },
    textures : {
        bar : {
            "1" : { x: 0, y: 0, width: 24, height: 16 },
            "0" : { x: 24, y: 0, width: 24, height: 16 },
        }
    }
});

Stage({
    name: 'lovebar',
    image: {
        src : './assets/lovebar.png',
    },
    textures : {
        bar : {
            "1" : { x: 0, y: 0, width: 24, height: 16 },
            "0" : { x: 24, y: 0, width: 24, height: 16 },
        }
    }
});

Stage(function (stage, display) {
    // create a WS connection to the server
    var socket = new WebSocket("ws://localhost:8081/ws/0");
    // handle deltas (changes to statistics from timers or interactions)
    function handleDelta(delta) {
        var type = delta.type;
        var amount = delta.amount;
        switch(type) {
            case "food":
                food_delta(amount);
                break;
            case "love":
                love_delta(amount);
                break;
        }
    }

    function handleHatch(hatch) {

    }

    socket.onmessage = function(event) {
      eventObj = JSON.parse(event.data);
      console.log(eventObj)
      switch(eventObj.type) {
          case "DELTA":
              handleDelta(eventObj.delta);
              break;
          case "HATCH":
              handleHatch(eventObj.hatch);
              break;
          case "MESSAGE":
              handleMessage(eventObj.message);
              break;
          case "REFRESH":
              handleRefresh(eventObj.state);
              break
      }
    };
    stage.viewbox(320, 240, mode='in-pad');
    stage.on('viewport', function(viewport) {
        console.log("resizing...");
    });
    stage.MAX_ELAPSE = 20;


    function handleMessage(message) {
        console.log(message);
    }

    function handleRefresh(state) {
        console.log(state);
        current_food = state.food;
        max_food = state.food_max;
        current_love = state.love;
        max_love = state.love_max;
        update_all()
    }

    function update_all() {
        update_foodbar()
        update_lovebar()
    }
    // foodbar functions
    function update_foodbar() {
        foodbar.setValue(createBarString(current_food, max_food));
    }
    var foodbar = Stage.string('foodbar:bar').appendTo(stage).pin({
        alignX : 0,
        alignY : 0,
        offsetX : 16,
        offsetY : 16,
    });
    foodbar.on('click', update_foodbar);
    update_foodbar();
    function food_delta(amount) {
        current_food += amount
        current_food = Math.min(Math.max(current_food, 0), max_food)
        update_foodbar()
    };

    // lovebar functions
    function update_lovebar() {
        lovebar.setValue(createBarString(current_love, max_love));
    }
    function love_delta(amount) {
        console.log("amount: ", amount);
        current_love += amount
        current_love = Math.min(Math.max(current_love, 0), max_love)
        console.log("current love: ", current_love, "max love: ", max_love)
        update_lovebar()
    };
    var lovebar = Stage.string('lovebar:bar').appendTo(stage).pin({
        alignX : 0,
        alignY : 0,
        offsetX : 16,
        offsetY : 40,
    });
    lovebar.on('click', update_lovebar);
    update_lovebar();
});
