#include <pthread.h>
#include <stdint.h>
#include <unistd.h>

#if defined(__linux__)
#include <sys/syscall.h>
#elif defined(__FreeBSD__) || defined(__OpenBSD__)
#include <pthread_np.h>
#elif defined(__NetBSD__)
#include <lwp.h>
#else
#endif

#include <roc/log.h>

#include "_cgo_export.h"

unsigned long long rocGoThreadID() {
#if defined(SYS_gettid)
    return (unsigned long long)(pid_t)syscall(SYS_gettid);
#elif defined(__FreeBSD__)
    return (unsigned long long)pthread_getthreadid_np();
#elif defined(__NetBSD__)
    return (unsigned long long)_lwp_self();
#elif defined(__APPLE__)
    uint64_t tid = 0;
    pthread_threadid_np(NULL, &tid);
    return (unsigned long long)tid;
#elif defined(__ANDROID__)
    return (unsigned long long)gettid();
#else
    return (unsigned long long)pthread_self();
#endif
}

void rocGoLogHandlerProxy(const roc_log_message* message, void* argument) {
    rocGoLogHandler((roc_log_message*)message);
}
