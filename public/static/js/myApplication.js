
function plural_form(number, after) {
    var cases = [2, 0, 1, 1, 1, 2];
    return number+' '+after[ (number%100>4 && number%100<20) ? 2: cases[Math.min(number%10, 5)] ];
}

angular.module('chat', [
    'ngAnimate',
    'ngWebSocket'
])
    .factory('Messages', function($websocket) {
        var ws = $websocket('ws://localhost:8070/realtime');
        var collection = [];
        var FirstWinner,SecondWinner,LastWinner = null;

        ws.onMessage(function(event) {
            var res;
            try {
                res = JSON.parse(event.data);
                console.log('message: ',Array.isArray(res.body.data),res);
                switch (res.action){
                    case "emitAll":
                        console.log('[emitAll]: ',res.body.data != null,res.body.data);
                        if(Array.isArray(res.body.data) && res.body.data.length > 0){
                            collection = [];
                            res.body.data.forEach(function (value) {
                                if (value.chat_id == 1){
                                    return;
                                }
                                collection.push(value);
                            })
                            console.log('[collection]',collection);
                        }
                        break;
                    case "sendUsers":
                        console.log('[sendUsers]: ',res.body.data != null,res.body.data);
                        if(Array.isArray(res.body.data) && res.body.data.length > 0){
                            collection = [];
                            res.body.data.forEach(function (value) {
                                if (value.chat_id == 1){
                                    return;
                                }
                                collection.push(value);
                            })
                            console.log('[collection]',collection);
                        }
                        break;

                    case "sendWebThreeWinner":
                        LastWinner = res.body.data;
                        console.log("LastWinner",LastWinner);
                        break;
                    case "sendWebTwoWinner":
                        SecondWinner = res.body.data;
                        console.log("SecondWinner",SecondWinner);
                        break;
                    case "sendWebOneWinner":
                        FirstWinner = res.body.data;
                        console.log("FirstWinner",FirstWinner);
                        break;
                    default:
                        console.log(res)
                }
            } catch(e) {
                res = {'username': 'anonymous', 'message': event.data};
            }
        });

        ws.onError(function(event) {
            console.log('connection Error', event);
        });

        ws.onClose(function(event) {
            ws = $websocket('ws://localhost:8070/realtime');
            console.log('connection closed', event);
        });

        ws.onOpen(function() {
            console.log('connection open');
            ws.send({
                action:"getUsers",
                body:{
                    message:"Дай мне всех пользователей",
                    data:{}
                },
                user:{
                    chat_id:1,
                    phone:"",
                    name:"",
                    login:"WebInterface",
                    photo:"/static/noavatar.png",
                    status:false
                }
            });
        });
        // setTimeout(function() {
        //   ws.close();
        // }, 500)

        return {
            collection: collection,
            SecondWinner: SecondWinner,
            FirstWinner: FirstWinner,
            LastWinner: LastWinner,
            status: function() {
                return ws.readyState;
            },
            getFirstWin: function () {
                return FirstWinner;
            },
            getLastWin: function () {
                return LastWinner;
            },
            getSecond: function () {
                return SecondWinner;
            },
            getConnects: function () {
                var c = collection.length;
                return 'В комнате : ' + plural_form(c,['пользователь','пользователя','пользователей']);
            },
            getUsers: function () {
                return collection;
            },
            getPhotoPath: function(photo) {
              if(typeof photo == "undefined"){
                  return "/static/noavatar.png"
              }
              return photo
            },
            getStatus: function (status) {
                if (status == 1) {
                    return 'glyphicon glyphicon-ok-sign text-success'
                }
                return 'glyphicon glyphicon-remove-sign text-danger'
            },
            send: function(message) {
                if (angular.isString(message)) {
                    ws.send(message);
                }
                else if (angular.isObject(message)) {
                    ws.send(JSON.stringify(message));
                }
            }

        };
    })
    .controller('MessengerController', function($scope, Messages) {
        console.log('MessengerController',Messages);
        $scope.Messages = Messages;

    })
    .filter('capitalize', function() {
        function capWord(txt) {
            console.log(txt)
            return txt.charAt(0).toUpperCase() + txt.substr(1).toLowerCase();
        }
        return function(input, isEveryWord) {
            return (!input) ? '' : (!isEveryWord) ? capWord(input) : input.replace(/([^\W_]+[^\s-]*) */g, capWord);
        };
    });