#include <roc/log.h>
#include "_cgo_export.h"

void rocGoLogHandlerProxy(const roc_log_message* message, void* argument) {
    rocGoLogHandler((roc_log_message*)message);
}