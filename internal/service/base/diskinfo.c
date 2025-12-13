#include <windows.h>
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <winioctl.h>

#include "diskinfo.h"

STORAGE_DEVICE_INFO* get_physicaldrive_info_struct_c(int driveIndex) {
    char devicePath[64];
    snprintf(devicePath, sizeof(devicePath), "\\\\.\\PhysicalDrive%d", driveIndex);

    HANDLE hDevice = CreateFileA(
        devicePath,
        GENERIC_READ,
        FILE_SHARE_READ | FILE_SHARE_WRITE,
        NULL,
        OPEN_EXISTING,
        0,
        NULL
    );

    if (hDevice == INVALID_HANDLE_VALUE) {
        return NULL;
    }

    STORAGE_DEVICE_INFO* info =
        (STORAGE_DEVICE_INFO*)calloc(1, sizeof(STORAGE_DEVICE_INFO));
    if (!info) {
        CloseHandle(hDevice);
        return NULL;
    }

    BYTE buffer[1024];
    DWORD bytesReturned = 0;

    STORAGE_PROPERTY_QUERY query;
    ZeroMemory(&query, sizeof(query));
    query.PropertyId = StorageAdapterProperty;
    query.QueryType = PropertyStandardQuery;

    if (DeviceIoControl(
        hDevice,
        IOCTL_STORAGE_QUERY_PROPERTY,
        &query,
        sizeof(query),
        buffer,
        sizeof(buffer),
        &bytesReturned,
        NULL
    )) {
        STORAGE_ADAPTER_DESCRIPTOR* adapter =
            (STORAGE_ADAPTER_DESCRIPTOR*)buffer;

        info->busType = adapter->BusType;
        info->commandQueueing = adapter->CommandQueueing;
        info->maxTransfer = adapter->MaximumTransferLength;
    }

    DISK_GEOMETRY geom;
    ZeroMemory(&geom, sizeof(geom));

    if (DeviceIoControl(
        hDevice,
        IOCTL_DISK_GET_DRIVE_GEOMETRY,
        NULL,
        0,
        &geom,
        sizeof(geom),
        &bytesReturned,
        NULL
    )) {
        info->bytesPerSector = geom.BytesPerSector;
    }

    CloseHandle(hDevice);
    return info;
}
