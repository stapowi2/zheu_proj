self.addEventListener('push', function(event) {
    const data = event.data ? event.data.text() : 'Новое уведомление!';
    event.waitUntil(
        self.registration.showNotification('Система ЖЭУ', {
            body: data,
        })
    );
});

