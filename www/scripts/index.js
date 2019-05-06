$('.form').find('input, textarea').on('keyup blur focus', function(e) {
    var $this = $(this),
        label = $this.prev('label');

    if (e.type === 'keyup') {
        if ($this.val() === '') {
            label.removeClass('active highlight');
        } else {
            label.addClass('active highlight');
        }
    } else if (e.type === 'blur') {
        if ($this.val() === '') {
            label.removeClass('active highlight');
        } else {
            label.removeClass('highlight');
        }
    } else if (e.type === 'focus') {
        if ($this.val() === '') {
            label.removeClass('highlight');
        } else if ($this.val() !== '') {
            label.addClass('highlight');
        }
    }
});

$('.tab a').on('click', function(e) {
    e.preventDefault();

    $(this).parent().addClass('active');
    $(this).parent().siblings().removeClass('active');

    target = $(this).attr('href');

    $('.tab-content > div').not(target).hide();

    $(target).fadeIn(600);
});

function generateCaptcha() {
  $.ajax({
    method: "POST",
    url: "/generate-captcha",
  }).done(function(captcha) {
    $("#captcha-img").attr("src", captcha);
  });
}

// my js code!!!
$("#sign-up-form").on("submit", (e) => {
  e.preventDefault();
  e.stopPropagation();

  var emailExists = false;

  $.when($.ajax({
    method: "POST",
    url: "/check-if-email-exists",
    data: "email=" + $("#emailSignUp").val().toLowerCase(),
  })
    .done(function(msg) {
      if (msg.length === 0) { return }
      emailExists = true;
      alert("User with such email already exists!");
      $("#emailSignUp").val("");
      $("#emailSignUp").focus();
    }))
      .then(function(){
        if (emailExists) { return }
        var bits = sjcl.hash.sha256.hash($("#passwordSignUp").val());
        var passwordHash = sjcl.codec.hex.fromBits(bits);

        $.ajax({
          method: "POST",
          url: "/sign-up",
          data: "firstName=" + $("#firstName").val() + "&lastName=" + $("#lastName").val() + "&email=" +
            $("#emailSignUp").val().toLowerCase() + "&password=" +  passwordHash
        })
          .done(function(msg) {
            if (msg.length > 0) {
              alert(msg);
            } else {
              window.location = "/registered";
            }
          });
      });
});

$("#sign-in-form").on("submit", (e) => {
  e.preventDefault();
  e.stopPropagation();

  var bits = sjcl.hash.sha256.hash($("#passwordSignIn").val());
  var passwordHash = sjcl.codec.hex.fromBits(bits);

  $.ajax({
    method: "POST",
    url: "/login",
    data: "email=" + $("#emailSignIn").val().toLowerCase() + "&password=" + passwordHash + "&captcha="
      + $("#captcha-input").val()
  }).done(function(msg) {
      console.log("zzz!!!");
      if (msg === "0" || msg === "1" || msg === "2") {
        if (msg === "0") {
          alert("Incorrect email or password. Try again!");
          $("#emailSignIn").val("");
          $("#passwordSignIn").val("");
          $("#emailSignIn").focus();
        } else if (msg === "1") {
          alert("Please confirm your email!");
        } else {
          alert("Your captcha solution is incorrect. Please try again!");
          generateCaptcha();
          $("#captcha-input").val("");
          $("#captcha-input").focus();
        }
      } else {
        window.location = "/contracts";
      }
    });
});

$("#login-tab").on("click", () => {
  generateCaptcha();
});


// update captcha every minute
var captchaTimer = setInterval(generateCaptcha, 60000);

$("#captcha-img").on("click", () => {
  generateCaptcha();
  clearInterval(captchaTimer);
  captchaTimer = setInterval(generateCaptcha, 60000);
});
