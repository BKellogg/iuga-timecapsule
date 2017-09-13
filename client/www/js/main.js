"use strict";

const CAPSULE_URL = 'https://' + 'iuga-tc.brendanjkellogg.com' + '/capsule';

var netidInput = $('#form_netid');
var gradDateInput = $('#form_graddate');
var messageInput = $('#form_message');
var tcForm = $('.tc-form');
var errorBox = $('.js-errorbox');

// prepopulate the date field with a value that is accurate for most people
gradDateInput.val('15 June, 2019');

$('.datepicker').pickadate({
    selectMonths: true, // Creates a dropdown to control month
    selectYears: 15, // Creates a dropdown of 15 years to control year,
    clear: 'Clear',
    close: 'Ok',
    closeOnSelect: true, // Close upon selecting a date
    min: '15 June, 2018',
    max: '30 June, 2020'
});

// process and submit the time capsule form
tcForm.submit((evt) => {
    evt.preventDefault();
    if (!errorBox.hasClass('hidden')) {
        errorBox.addClass('hidden');
    }
    var date = moment(gradDateInput.val(), 'D MMMM, YYYY').format('YYYY/MM/DD')
    tcForm.slideToggle();
    var paylaod = {
        'netID': netidInput.val(),
        'gradDate': date,
        'message': messageInput.val()
    }
    var reqHeaders = new Headers();
    reqHeaders.append('Content-Type', 'application/json');
    fetch(CAPSULE_URL, {
        method: 'POST',
        headers: reqHeaders,
        body: JSON.stringify(paylaod)
    })
    .then((res) => {
        if (res.ok) {
            return res.json();
        }
        return res.text().then((t) => Promise.reject(t));
    })
    .then((data) => {
        handleSuccess(data);
    })
    .catch((err) => {
        console.warn(err);
        errorBox.text(err);
        errorBox.removeClass('hidden');
        tcForm.slideToggle();
    });
});

// handleSuccess displays the success box with a confirmation message
function handleSuccess(data) {
    $('.js-successbox').removeClass('hidden');
    $('.success-date').text(moment(data.gradDate, 'YYYY/MM/DD').format('MMMM Do, YYYY'));
    $('.success-email').text(data.email);
}