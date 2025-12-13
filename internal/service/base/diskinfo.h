#pragma once

typedef struct {
    int busType;
    int commandQueueing;
    unsigned long maxTransfer;
    unsigned long bytesPerSector;
} STORAGE_DEVICE_INFO;

STORAGE_DEVICE_INFO* get_physicaldrive_info_struct_c(int driveIndex);
