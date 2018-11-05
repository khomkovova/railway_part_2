



new Vue({
    el: '#app',
    data:{
        error:"",
        status:"",
        comments:""
    },
    methods:{

        submit(){
            // username = this.username;
            // password = this.password;
            // (function($){

            console.log("asdfasdf");
            axios({ method: "POST", "url": "http:api/addcomments", data:'{"comments":'+ '"' + this.comments + '"}'}).then(result => {
                // this.jsonInfo = result.data;
                // var obj = JSON.parse(result.data);
                // console.error(result.data);
                // window.location.assign('/updatefirmware');
                this.status = "Your comment send";
                // console.error(result.data);
            }, error => {
                this.error="Bad username or password";
                // console.error(error);
            });

        }
    }
});