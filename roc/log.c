#include <roc/log.h>
#include "_cgo_export.h"

void rocGoLogHandlerProxy(roc_log_level level, char* component, char* message) {
    rocGoLogHandler(level, component, message);
}
