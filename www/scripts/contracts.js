// misha froloff's code
function openDivs(numb){
    $(".main_container").hide()
    $(".jumb").hide()
    $("#main_container_" + numb).show()
    $("#jumb_" + numb).show()
    if (numb != 0){
      $("#aside").hide()
    } else{
      $("#aside").show()
    }
} openDivs(4)

function loadAttrs(type){
    $(".smart_attrs").hide()
    $("#" + type + "_attrs").show()
    $("#simpleTransactionName").html(type + " transaction")
} loadAttrs("Deferred")

function addReceiver(button) {
    $('<div class="form-row receiver-additional"><div class="form-group col-md-10"><input type="text" class="form-control masReceiverValue" placeholder="Additional public key" pattern="[a-fA-F0-9]{66}" title="Please type the 66 digit HEX code." required></div><div class="form-group"><button type="button" onClick = "delPrev(this);" style = "margin-left:3px;" class="btn btn-danger">Delete</button></div></div>').insertBefore($(button))
}
function addReceiverWithValue(button, value) {
    $('<div class="form-row receiver-additional"><div class="form-group col-md-10"><input type="text" class="form-control masReceiverValue" placeholder="Additional public key" value="' + value + '"></div><div class="form-group"><button type="button" onClick = "delPrev(this);" style = "margin-left:3px;" class="btn btn-danger">Delete</button></div></div>').insertBefore($(button))
}
function delPrev(button){
    $(button).parent().parent().remove()
}

$(function () {
    $('#datetimepicker, #datetimepicker2, #datetimepicker3, #datetimepicker4').datetimepicker({
        inline: true,
        sideBySide: true
    });
});

// ENABLE ACE EDITOR
var aceEditor = ace.edit("editor");
aceEditor.setTheme("ace/theme/pastel_on_dark");
aceEditor.getSession().setMode("ace/mode/python");


// *************
// my js code!!!
// *************

// hide the simple smart-contract edit label
$("#simple-smc-edit-label").hide();

// check if picked date is more than today
function CheckIsContractDateActual(date) {
  today = new Date();
  if (date >= today) {
    return true
  }
  alert("Please choose date in the future!");
  $('#datetimepicker').data("datetimepicker").setDate(today);
  return false;
}

$("#deferred-payment-form").on("submit", (event) => {
  event.preventDefault();
  event.stopPropagation();

  // get date from datetimepicker
  var date = $("#datetimepicker").data("datetimepicker").getDate();
  // check contract date
  if (CheckIsContractDateActual(date) == false) {
    return;
  }

  var formatted = date.getFullYear() + "-" + ("0" + (date.getMonth() + 1)).slice(-2) + "-" +
    ("0" + date.getDate()).slice(-2) + " " + ("0" + date.getHours()).slice(-2) + ":" +
    ("0" + date.getMinutes()).slice(-2) + ":00";
  console.log(formatted);

  id = "not_exist";
  if ($("#data").attr("smart-contract-ID") !== "") {
    id = $("#data").attr("smart-contract-ID");
  }

  $.ajax({
    method: "POST",
    url: "/generate-contract",
    data: {ContractType: "deferred", ID: id, ContractDate: formatted, Receiver: $("#inputReceiver").val(),
      Data: $("#inputText").val(), TransactionMessage: $("#inputTransaction").val(),}
  }).done(function(msg) {
    if (msg === "Error") {
      alert("Server error!");
    } else {
      if (msg === "ok") {
        window.location = "/contracts";
      } else {
        document.write(msg);
      }
    }
  });
});

$("#condition-payment-form").on("submit", (event) => {
  event.preventDefault();
  event.stopPropagation();

  var date = $("#datetimepicker2").data("datetimepicker").getDate();
  // check contract date
  if (CheckIsContractDateActual(date) == false) {
    return;
  }

  var formatted = date.getFullYear() + "-" + ("0" + (date.getMonth() + 1)).slice(-2) + "-" +
    ("0" + date.getDate()).slice(-2) + " " + ("0" + date.getHours()).slice(-2) + ":" +
    ("0" + date.getMinutes()).slice(-2) + ":00";
  console.log(formatted);

  id = "not_exist";
  if ($("#data").attr("smart-contract-ID") !== "") {
    id = $("#data").attr("smart-contract-ID");
  }

  // check select value (can not be less than zero)
  selectValueInt = parseInt($("#selectValue").val(), 10);
  if (selectValueInt <= 0) {
    alert("Please type a positive number!");
    $("#selectValue").val("");
    return;
  }

  $.ajax({
    method: "POST",
    url: "/generate-contract",
    data: {ContractType: "condition", ID: id, ContractDate: formatted, SelectCondition: $("#selectCondition").val(),
      SelectOperator: $("#selectOperator").val(), SelectValue: $("#selectValue").val(),
      Receiver: $("#inputReceiver2").val(), Data: $("#inputText2").val(),
      TransactionMessage: $("#inputTransaction2").val(),}
  }).done(function(msg) {
    if (msg === "Error") {
      alert("Server error!");
    } else {
      if (msg === "ok") {
        window.location = "/contracts";
      } else {
        document.write(msg);
      }
    }
  });
});

$("#auto-payment-form").on("submit", (event) => {
  event.preventDefault();
  event.stopPropagation();

  var date = $("#datetimepicker3").data("datetimepicker").getDate();
  // check contract date
  if (CheckIsContractDateActual(date) == false) {
    return;
  }

  var formatted = date.getFullYear() + "-" + ("0" + (date.getMonth() + 1)).slice(-2) + "-" +
    ("0" + date.getDate()).slice(-2) + " " + ("0" + date.getHours()).slice(-2) + ":" +
    ("0" + date.getMinutes()).slice(-2) + ":00";
  console.log(formatted);

  id = "not_exist";
  if ($("#data").attr("smart-contract-ID") !== "") {
    id = $("#data").attr("smart-contract-ID");
  }

  $.ajax({
    method: "POST",
    url: "/generate-contract",
    data: {ContractType: "auto", ID: id, ContractDate: formatted, AutoPaymentMode: $("#selectSettings").val(),
    Receiver: $("#inputReceiver3").val(), Data: $("#inputText3").val(),
    TransactionMessage: $("#inputTransaction3").val(),}
  }).done(function(msg) {
    if (msg === "Error") {
      alert("Server error!");
    } else {
      if (msg === "ok") {
        window.location = "/contracts";
      } else {
        document.write(msg);
      }
    }
  });
});

$("#collective-payment-form").on("submit", (event) => {
  event.preventDefault();
  event.stopPropagation();

  var date = $("#datetimepicker4").data("datetimepicker").getDate();
  // check contract date
  if (CheckIsContractDateActual(date) == false) {
    return;
  }

  var formatted = date.getFullYear() + "-" + ("0" + (date.getMonth() + 1)).slice(-2) + "-" +
    ("0" + date.getDate()).slice(-2) + " " + ("0" + date.getHours()).slice(-2) + ":" +
    ("0" + date.getMinutes()).slice(-2) + ":00";
  console.log(formatted);

  var receivers = "";
  var allFieldsAreNotEmpty = true;
  $(".masReceiverValue").each(
    function(index) {
      if ($(this).val() === "") {
        allFieldsAreNotEmpty = false;
        return;
      }
      receivers += "'" + $(this).val() + "',";
    }
  );
  if (allFieldsAreNotEmpty === false) {
    alert("Please fill all receiver fields!");
    return;
  }
  receivers = "[" + receivers.slice(0, -1) + "]";

  id = "not_exist";
  if ($("#data").attr("smart-contract-ID") !== "") {
    id = $("#data").attr("smart-contract-ID");
  }

  $.ajax({
    method: "POST",
    url: "/generate-contract",
    data: {ContractType: "collective", ID: id, ContractDate: formatted, Receivers: receivers,
    Data: $("#inputText4").val(), TransactionMessage: $("#inputTransaction4").val(),}
  }).done(function(msg) {
    if (msg === "Error") {
      alert("Server error!");
    } else {
      if (msg === "ok") {
        window.location = "/contracts";
      } else {
        document.write(msg);
      }
    }
  });
});

function UpdateStatusOfSmartContract(id, statusNew) {
  $.ajax({
    method: "POST",
    url: "/updateStatusOfSmartContract",
    data: {ID: id, Status: statusNew},
  }).done(function(msg) {
    if (msg === "Error") {
      alert("Server error!");
    }
  });
}

function pythonToJSReceivers(receivers) {
  bracketsRemoved = receivers.slice(1, -1);
  receiversArray = bracketsRemoved.split(",");
  for (var i=0; i<receiversArray.length; i++) {
    receiversArray[i] = receiversArray[i].slice(1, -1);
  }
  return receiversArray
}

var aceEditorChanged = false;

function BuildSmartContractsTable() {
  // console.log("'my-smart-contracts tab opened!!!'");
  $.ajax({
    method: "POST",
    url: "/get-all-smart-contracts",
  }).done(function(data) {
    console.log("data=[" + data + "]")
    if (data === "Error") {
      alert("Server error!");
    } else if (data === "[]") {
      $("#my-smart-contracts-table").html("no smart contracts added yet");
    } else if (data === "please log in to the system") {
      $("#my-smart-contracts-table").html(`please <a href="/">log in</a> to the system`);
      // $("#my-smart-contracts-table").html("please log in to the system");
    } else {
      // populate 'my-smart-contracts' table here
      dataParsed = $.parseJSON(data);
      $("#my-smart-contracts-table > tbody").html("");
      $.each(dataParsed, function(i, smartContract) {
        $('<tr valign="center">').append(
          '<td>' + i + '</td>' +
          '<td id="smc_id_' + i + '">' + smartContract.ID + '</td>' +
          '<td id="smc_status_' + i + '">' + smartContract.Status + '</td>' +
          '<td>' + smartContract.Type + ' payment</td>' +
          '<td>' + smartContract.CreationDate + '</td>' +
          '<td>' + smartContract.Price + ' NZT</td>' +
          '<td>' + smartContract.LastStarted + '</td>' +
          '<td>' + smartContract.Comment + '</td>' +
          `<td>
            <div class="mbtn-group">
              <div class = "row"><button class = "btn btn-sm mbtn" id="smc_start_` + i + `">Start</button></div>
              <div class = "row"><button class = "btn btn-sm mbtn" id="smc_pause_` + i + `">Pause</button></div>
              <div class = "row"><button class = "btn btn-sm mbtn" id="smc_edit_` + i + `">Edit</button></div>
              <div class = "row"><button class = "btn btn-sm mbtn" id="smc_delete_` + i + `">Delete</button></div>
            </div>
          </td>`
        ).appendTo('#my-smart-contracts-table');

        // implement smart-contract start button handler
        $("#smc_start_" + i).on("click", () => {
          console.log("smart-contract #" + i + " is going to start!");
          $.ajax({
            method: "POST",
            url: "/readSmartContract",
            data: "ID=" + $('#smc_id_' + i).text(),
          }).done(function(transaction) {
              // get smart-contract price
              transactionParsed = JSON.parse(transaction);
              tran = transactionParsed.transaction;
              smcStatus = transactionParsed.smcstatus;
              smcWasStarted = transactionParsed.smcwasstarted;

              tranParsed = JSON.parse(tran);
              smcCode = tranParsed.CODE;
              smcTimestamp = tranParsed.TST;

              // private key is strongly specified; testing only!!!
              var privateKey = 'a376bdbcca0cda0ed27c7f5a434444276b2379b25f504c3ba9c3f4f52a6a5d1c';
              var publicKey = secp256k1.getPublicKey(privateKey);

              if (smcStatus === "READY TO START" && smcWasStarted === false) {
                // remove last character (curly bracket) and add public key as sender
                // '-2' because sometimes slice(0, -1) doesn't delete the last curly bracket of json
                transaction = tran.slice(0, -2) + `,"SENDER":"` + publicKey +  `"}`;

                // sign transaction
                var signatureHex = secp256k1.sign(transaction, privateKey);

                // remove the last character (curly bracket) and add sign
                transaction = transaction.slice(0, -1) + `,"SIGNATURE":"` + signatureHex + `"}`;

                // send smart-contract transaction to PVM
                $.ajax({
                  method: "POST",
                  url: "/sendTransactionToPVM",
                  data: {Transaction: transaction, ID: $('#smc_id_' + i).text(), SetWasStarted: 1,},
                }).done(function(result) {
                  if (result === "error") {
                    alert("Error!");
                  } else {
                    UpdateStatusOfSmartContract($('#smc_id_' + i).text(), "WORKING");
                    $("#smc_status_" + i).text("WORKING");
                    alert("Transaction was sent to PVM successfully!");
                  }
                });
              } else if (smcStatus === "PAUSED" && smcWasStarted === true) {
                smcCodeSignature = secp256k1.sign(smcCode, privateKey);
                actionJSON = `{"SENDER":"` + publicKey + `","ACTION":"START","SMS":"` + smcCodeSignature + `","TST":"`
                  + smcTimestamp + `"}`;

                actionJSONSignature = secp256k1.sign(actionJSON, privateKey);
                actionJSON = actionJSON.slice(0, -1) + `,"SIGNATURE":"` + actionJSONSignature + `"}`;

                // send smart-contract action to PVM
                $.ajax({
                  method: "POST",
                  url: "/sendActionToPVM",
                  data: {"ActionJSON": actionJSON,},
                }).done(function(result) {
                  if (result === "error") {
                    alert("Error!");
                  } else {
                    UpdateStatusOfSmartContract($('#smc_id_' + i).text(), "WORKING");
                    $("#smc_status_" + i).text("WORKING");
                    alert("Action was sent to PVM successfully!");
                  }
                });
              } else if (smcStatus === "WORKING" && smcWasStarted === true) {
                alert("Smart-contract is already started!");
              } else {
                alert("Unknown case. Please contact to administrator.");
              }
          });
        });

        // implement smart-contract delete button handler
        $("#smc_delete_" + i).on("click", () => {
          console.log("smart-contract #" + i + " is going to delete!");

          if (!confirm("Are you sure?")) {
            return
          }

          $.ajax({
            method: "POST",
            url: "/deleteSmartContract",
            data: "ID=" + $('#smc_id_' + i).text(),
          }).done(function(msg) {
            if (msg === "Error") {
              alert("Server error!");
            } else {
              // $("#my-smart-contracts-table").empty();
              $("#my-smart-contracts-table > tbody").html("");
              BuildSmartContractsTable();
            }
          });
        });

        // implement smart-contract edit button handler
        $("#smc_edit_" + i).on("click", () => {
          console.log("smart-contract #" + i + " is going to edit!");

          $.ajax({
            method: "POST",
            url: "/editSmartContract",
            data: "ID=" + $('#smc_id_' + i).text(),
          }).done(function(smartContract) {
            if (smartContract === "Error") {
              alert("Server error!");
            } else {
              // parse json from server
              smartContractParsed = $.parseJSON(smartContract);

              // set id attr
              $("#data").attr("smart-contract-ID", smartContractParsed.ID);

              // show label with smc ID
              $("#simple-smc-edit-label").text("Edit the simple smart-contract with ID=" + smartContractParsed.ID);
              $("#simple-smc-edit-label").show();

              if (smartContractParsed.Type === "custom") {
                // update the title of the tab
                $("#titleOfTheAceEditor").text("Edit the smart-contract with ID=" + smartContractParsed.ID + ":");
                // set the value of ace editor
                aceEditor.setValue(smartContractParsed.Code, -1);
                // reset aceEditorChanged
                aceEditorChanged = false;
                // hide pylint errors block
                HideSmartContractErrors();
                // open the tab with ace editor
                openDivs(2);
              } else if (smartContractParsed.Type === "deferred") {
                smartContractDate = new Date(smartContractParsed.Data.ContractDate.replace(" ", "T"));
                console.log("date=" + smartContractParsed.Data.ContractDate.replace(" ", "T"));
                $('#datetimepicker').data("datetimepicker").setDate(smartContractDate);

                // $("#inputTokenCount").val(smartContractParsed.Data.TokenCount);
                // $("#inputTokenType").val(smartContractParsed.Data.TokenType);
                $("#inputReceiver").val(smartContractParsed.Data.Receiver);
                $("#inputText").val(smartContractParsed.Data.Data);
                $("#inputTransaction").val(smartContractParsed.Data.TransactionMessage);

                // open tab with simple editor
                openDivs(1);
                loadAttrs('Deferred');
              } else if (smartContractParsed.Type === "condition") {
                smartContractDate = new Date(smartContractParsed.Data.ContractDate.replace(" ", "T"));
                $('#datetimepicker2').data("datetimepicker").setDate(smartContractDate);

                $("#inputTokenCount2").val(smartContractParsed.Data.TokenCount);
                $("#inputTokenType2").val(smartContractParsed.Data.TokenType);
                $("#selectCondition").val(smartContractParsed.Data.SelectCondition);
                $("#selectOperator").val(smartContractParsed.Data.SelectOperator);
                $("#selectValue").val(smartContractParsed.Data.SelectValue);
                $("#inputReceiver2").val(smartContractParsed.Data.Receiver);
                $("#inputText2").val(smartContractParsed.Data.TransactionMessage);

                // open tab with simple editor
                openDivs(1);
                loadAttrs('Condition');
              } else if (smartContractParsed.Type === "auto") {
                smartContractDate = new Date(smartContractParsed.Data.ContractDate.replace(" ", "T"));
                $("#datetimepicker3").data("datetimepicker").setDate(smartContractDate);

                $("#selectSettings").val(smartContractParsed.Data.AutoPaymentMode);
                $("#inputTokenCount3").val(smartContractParsed.Data.TokenCount);
                $("#inputTokenType3").val(smartContractParsed.Data.TokenType);
                $("#inputReceiver3").val(smartContractParsed.Data.Receiver);
                $("#inputText3").val(smartContractParsed.Data.TransactionMessage);

                // open tab with simple editor
                openDivs(1);
                loadAttrs('Auto');
              } else if (smartContractParsed.Type === "collective") {
                // remove additional receiver inputs
                $(".receiver-additional").remove();

                smartContractDate = new Date(smartContractParsed.Data.ContractDate.replace(" ", "T"));
                $("#datetimepicker4").data("datetimepicker").setDate(smartContractDate);
a
                $("#inputTokenCount4").val(smartContractParsed.Data.TokenCount);
                $("#inputTokenType4").val(smartContractParsed.Data.TokenType);

                // insert receivers here...
                jsReceivers = pythonToJSReceivers(smartContractParsed.Data.Receivers);
                $("#collective-receiver-first").val(jsReceivers[0]);
                for (var i=1; i < jsReceivers.length; i++) {
                  addReceiverWithValue($("#add-receiver-btn"), jsReceivers[i]);
                }

                $("#inputText4").val(smartContractParsed.Data.TransactionMessage);

                // open tab with simple editor
                openDivs(1);
                loadAttrs('Collective');
              }
            }
          });
        });

        // implement smart-contract pause button handler
        $("#smc_pause_" + i).on("click", () => {
          console.log("smart-contract #" + i + " is going to pause!");
          $.ajax({
            method: "POST",
            url: "/readSmartContract",
            data: "ID=" + $('#smc_id_' + i).text(),
          }).done(function(transaction) {
            // get smart-contract price
            transactionParsed = JSON.parse(transaction);
            tran = transactionParsed.transaction;
            smcStatus = transactionParsed.smcstatus;
            console.log("smcStatus=" + smcStatus);
            smcWasStarted = transactionParsed.smcwasstarted;
            console.log("smcWasStarted=" + smcWasStarted);

            tranParsed = JSON.parse(tran);
            smcCode = tranParsed.CODE;
            smcTimestamp = tranParsed.TST;

            // private key is strongly specified; testing only!!!
            var privateKey = 'a376bdbcca0cda0ed27c7f5a434444276b2379b25f504c3ba9c3f4f52a6a5d1c';
            var publicKey = secp256k1.getPublicKey(privateKey);

            if (smcStatus === "READY TO START" && smcWasStarted === false) {
              alert("Smart-contract was not started yet. Please click 'start' to launch smart-contract!");
            } else if (smcStatus === "PAUSED" && smcWasStarted === true) {
              alert("Smart-contract is already paused.");
            } else if (smcStatus === "WORKING" && smcWasStarted === true) {
              smcCodeSignature = secp256k1.sign(smcCode, privateKey);
              actionJSON = `{"SENDER":"` + publicKey + `","ACTION":"PAUSE","SMS":"` + smcCodeSignature + `","TST":"`
                + smcTimestamp + `"}`;

              actionJSONSignature = secp256k1.sign(actionJSON, privateKey);
              actionJSON = actionJSON.slice(0, -1) + `,"SIGNATURE":"` + actionJSONSignature + `"}`;

              // send smart-contract action to PVM
              $.ajax({
                method: "POST",
                url: "/sendActionToPVM",
                data: {"ActionJSON": actionJSON,},
              }).done(function(result) {
                if (result === "error") {
                  alert("Error!");
                } else {
                  UpdateStatusOfSmartContract($('#smc_id_' + i).text(), "PAUSED");
                  $("#smc_status_" + i).text("PAUSED");
                  alert("Action was sent to PVM successfully!");
                }
              });
            } else {
              alert("Unknown case. Please contact to administrator.");
            }
          });
        });
      });
    }
  });
}

$("#my-smart-contracts-href").on("click", BuildSmartContractsTable);


// *** autosaving block ***

//setup before functions
var typingTimer;                //timer identifier
var doneTypingInterval = 50000;  //time in ms, 5 second for example

//user is "finished typing," do something
function doneTyping() {
  saveSmartContractIfChanged("autosaving_mode");
}

// declare event handlers
function aceEditorKeyUpHandler() {
  clearTimeout(typingTimer);
  typingTimer = setTimeout(doneTyping, doneTypingInterval);
}

function aceEditorKeyDownHandler() {
  clearTimeout(typingTimer);
}

$("#autosaving-mode").on("change", () => {
  if ($('#autosaving-mode').is(":checked")) {
    //on keyup, start the countdown
    aceEditor.textInput.getElement().addEventListener("keyup", aceEditorKeyUpHandler);
    //on keydown, clear the countdown
    aceEditor.textInput.getElement().addEventListener("keydown", aceEditorKeyDownHandler);
  } else {
    aceEditor.textInput.getElement().removeEventListener("keyup", aceEditorKeyUpHandler);
    aceEditor.textInput.getElement().removeEventListener("keydown", aceEditorKeyDownHandler);
  }
});


// smart-contracts error handling
$(".smart-contract-errors").hide();

var aceEditorMarkers = [];

function HideSmartContractErrors() {
  $(".smart-contract-errors").hide();
  for (i=0; i<aceEditorMarkers.length; i++) {
    aceEditor.session.removeMarker(aceEditorMarkers[i]);
  }
  aceEditorMarkers = [];
}

function GetSmartContractErrorsAndShow() {
  $.ajax({
    method: "POST",
    url: "/run-pylint",
    data: {SmartContractCode: aceEditor.getValue()},
  }).done(function(pylintErrors) {
    if (pylintErrors === "Error") {
      alert("Server error!");
    }
    // prepare errors report
    pylintErrorsParsed = $.parseJSON(pylintErrors);
    textErrors = "";

    $.each(pylintErrorsParsed, function(i, pylintError) {
      textErrors += pylintError.FullText + "\n";
      // highlight the error
      var Range = ace.require('ace/range').Range;
      marker = aceEditor.session.addMarker(new Range(pylintError.ErrorLineNumber - 1, 0, pylintError.ErrorLineNumber - 1, 1), "myMarker", "fullLine");
      aceEditorMarkers.push(marker);
    });

    if (textErrors === "") {
      alert("Your smart-contract isn't contain errors - good work!\n( :");
      return;
    }
    $("#pylint-errors").text(textErrors);
    $(".smart-contract-errors").show();
  });
}

function CheckSmartContractErrors() {
  if ($(".smart-contract-errors").is(":visible")) {
    HideSmartContractErrors();
  } else {
    GetSmartContractErrorsAndShow();
  }
}

$("#check-code-btn").on("click", () => {
  CheckSmartContractErrors();
});


// *******************************
// implement smart-contract saving

aceEditor.getSession().on("change", function() {
  aceEditorChanged = true;
});

function UpdateSmartContract(id, status, code, isAlert) {
  console.log("smart-contract with ID=" + id + " was saved (updated) successfully!");
  $.ajax({
    method: "POST",
    url: "/updateSmartContract",
    data: {ID: id, Status: status, SmartContractCode: code},
  }).done(function(msg) {
    if (msg === "Error") {
      alert("Server error!");
    } else {
      aceEditorChanged = false;
      if (isAlert) {
        alert("Smart-contract with ID=" + id + " was saved (updated) successfully!");
        window.location = "/contracts";
      }
    }
  });
}

function GenerateUniqueIDAndUpdateSmartContract(mode) {
  var smartContractID = -1;
  if ($("#data").attr("smart-contract-ID") === "") {
    $.ajax({
      method: "POST",
      url: "/generate-unique-id",
    }).done(function(uniqueID) {
      if (uniqueID === "Error") {
        alert("Server error!");
      }
      smartContractID = uniqueID;
    }).then(function() {
      if (mode === "click_event_mode" && !confirm("Are you sure?")) {
        return;
      }
      $("#data").attr("smart-contract-ID", smartContractID);
      $("#titleOfTheAceEditor").text("Edit the smart-contract with ID=" + smartContractID + ":");
      UpdateSmartContract(smartContractID, "DRAFT", aceEditor.getValue(), mode === "click_event_mode");
    });
  } else {
    if (mode === "click_event_mode" && !confirm("Are you sure?")) {
      return;
    }
    smartContractID = $("#data").attr("smart-contract-ID");
    UpdateSmartContract(smartContractID, "READY TO START", aceEditor.getValue(), mode === "click_event_mode");
  }
}

function saveSmartContractIfChanged(mode) {
  if (mode === "click_event_mode") {
    if (!aceEditorChanged && $("#data").attr("smart-contract-ID") !== "") {
      alert("Please make some changes before update.");
      return;
    } else if ($(".smart-contract-errors").is(":visible")) {
      alert("Please fix errors before save (update)!");
      return;
    }
  } else if (mode === "autosaving_mode") {
    // save smart-contract as draft (no error handling)
    GenerateUniqueIDAndUpdateSmartContract(mode);
    return;
  }

  $.ajax({
    method: "POST",
    url: "/run-pylint",
    data: {SmartContractCode: aceEditor.getValue()},
  }).done(function(pylintErrors) {
    if (pylintErrors === "Error") {
      alert("Server error!");
    }
    // prepare errors report
    pylintErrorsParsed = $.parseJSON(pylintErrors);
    textErrors = "";

    $.each(pylintErrorsParsed, function(i, pylintError) {
      textErrors += pylintError.FullText + "\n";
      // highlight the error
      var Range = ace.require('ace/range').Range;
      marker = aceEditor.session.addMarker(new Range(pylintError.ErrorLineNumber - 1, 0, pylintError.ErrorLineNumber - 1, 1), "myMarker", "fullLine");
      aceEditorMarkers.push(marker);
    });

    if (textErrors === "") {
      HideSmartContractErrors();
    } else {
      $("#pylint-errors").text(textErrors);
      $(".smart-contract-errors").show();
    }
  }).then(function() {
    if ($(".smart-contract-errors").is(":visible")) {
      alert("Please fix errors before save (update)!");
      return;
    }
    GenerateUniqueIDAndUpdateSmartContract(mode);
    aceEditorChanged = false;
  });
}

$("#saveSmartContractAce").on("click", () => {
  saveSmartContractIfChanged("click_event_mode");
});
