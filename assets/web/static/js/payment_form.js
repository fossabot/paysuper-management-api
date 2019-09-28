const PM_ACTIVE_CLASS = 'active';

$(function() {
    let $pmItem = $('div.payment-methods div.item');
    let $amountContainer = $('div.details div.amount div.main');
    let $commissionsContainer = $('div.details div.amount div.commissions');
    let $selectedPmInput = $('form#order-form input[name=payment_method_id]');

    $pmItem.each(function (i, el) {
        if (i === 0) {
            $('#' + $(el).data('form')).show();
        } else {
            $('#' + $(el).data('form')).find('input').attr({disabled: 'disabled'});
        }
    });

    if ($.trim($amountContainer.html()).length <= 0) {
        let amount = $pmItem.first().data('amount');
        let currency = $pmItem.first().data('currency');
        let vat = $pmItem.first().data('vat');
        let commission = $pmItem.first().data('commission');

        $amountContainer.html(amount + ' ' + currency);

        if (vat || commission) {
            $commissionsContainer.html('');
            $commissionsContainer.append('<div style="margin-bottom: 7px;">of them:</div>');

            if (vat) {
                $commissionsContainer.append('VAT ' + vat + ' ' + currency + '<br />');
            }

            if (commission) {
                $commissionsContainer.append('Commission ' + commission + ' ' + currency);
            }
        }
    }

    if ($selectedPmInput.val().length <= 0) {
        $selectedPmInput.val($pmItem.first().data('identifier'));
    }

    $pmItem.on('click', function () {
        $pmItem.removeClass(PM_ACTIVE_CLASS);
        $(this).addClass(PM_ACTIVE_CLASS);

        let $hideEl = $('div.payment-method-requisites').find('div.form');

        $hideEl.hide();
        $hideEl.find('input').attr({disabled: 'disabled'});

        let $showEl = $('#' + $(this).data('form'));

        $showEl.show();
        $showEl.find('input').removeAttr('disabled');

        let amount = $(this).data('amount');
        let currency = $(this).data('currency');
        let vat = $(this).data('vat');
        let commission = $(this).data('commission');

        $amountContainer.html(amount + ' ' + currency);

        if (vat || commission) {
            $commissionsContainer.html('');
            $commissionsContainer.append('<div style="margin-bottom: 7px;">of them:</div>');

            if (vat) {
                $commissionsContainer.append('VAT ' + vat + ' ' + currency + '<br />');
            }

            if (commission) {
                $commissionsContainer.append('Commission ' + commission + ' ' + currency);
            }
        }

        $selectedPmInput.val($(this).data('identifier'));
    });

    $('input.number').on('keyup', function () {
        $(this).val($(this).val().replace(/\D/g,''));
    });

    $('form#order-form').on('submit',function (e) {
        e.preventDefault();

        let data = new FormData($(this)[0]);
        let object = {};

        data.forEach(function(value, key){
            object[key] = value;
        });

        let $sForm = $('#redirect-form');

        $.ajax({
            url: '/api/v1/payment',
            type: 'post',
            data: JSON.stringify(object),
            cache: false,
            async: false,
            success: function(data) {
                if (!data.hasOwnProperty('redirect_url')) {
                    alert('Process will be stop, because we don\'t find required params.');
                    return;
                }

                $sForm.attr({action: data['redirect_url']}).submit();

                let centrifuge = new Centrifuge('wss://cf.tst.protocol.one/connection/websocket', {debug: true});
                centrifuge.setToken(token);

                const chanel = "payment:notify#"+object['order_id'];

                centrifuge.subscribe(chanel, function (message) {
                    alert("payment complete with status: " + message.data.status);
                });

                centrifuge.connect();
            },
            error: function(xhr) {
                let message = JSON.parse(xhr['responseText']);
                alert(message['error']);
            },
            contentType: 'application/json; charset=utf-8',
            processData: false,
            dataType: 'json'
        });
    });
});