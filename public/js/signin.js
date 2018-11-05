new Vue({
    el: '#app',
    data:{
        error:"",
        username:"",
        password:""
    },
    methods:{

        signin(){
            // username = this.username;
            // password = this.password;
            // (function($){
            console.log("asdfasdf");
            axios({ method: "POST", "url": "http:/api/signin", data:'{"username":"' + this.username + '", "password":"' + this.password + '"}' , withCredentials: true}).then(result => {
                // this.jsonInfo = result.data;
                // var obj = JSON.parse(result.data);
                // console.error(result.data);
                window.location.assign('/updatefirmware');
                console.error(result.data);
            }, error => {
                this.error="Bad username or password";
                // console.error(error);
            });

        }
    }
});