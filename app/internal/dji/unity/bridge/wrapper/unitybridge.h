#ifndef WRAPPER_UNITY_BRIDGE_H
#define WRAPPER_UNITY_BRIDGE_H

// Prototype for exported functions in unitybridge.dll.
//
// These are inferred from how C# code accesses them and there might be some
// unknown subtleties.
//
// As this is C code, boolean values are represented as ints and utf16 strings
// as unsigned shorts.
//
// TODO(bga): As of now, it mostly seems to work but it is not being
// extensivelly used. Basically I know it is not crashing but I can't know for
// sure yet that it is also doing what it is supposed to do correctly. 

// Prototype for event callback functions.
typedef void(*UnityEventCallbackFunc)(unsigned long long e, void* info,
                int length, unsigned long long tag);

// Unity Bridge construction and destruction.
void CreateUnityBridge(const char* name, int debuggable);
void DestroyUnityBridge();

// Unity Bridge initialization and uninitialization. The typo in the
// uninitialization function is present is the symbol exported.
int UnityBridgeInitialize();
void UnityBridgeUninitialze();

// Sets the callback function for specific events.
void UnitySetEventCallback(unsigned long long e,
                UnityEventCallbackFunc callback);

// Sends events that might be routed to the Robomaster S1. In the first function
// the data pointed by the info pointer might be changed during the call (if it
// is not NULL, that is). The other two are just wrappers for easier sending of
// strings and numbers (as int64, a.k.a. unsigned long long).
void UnitySendEvent(unsigned long long e, void* info, unsigned long long tag);
void UnitySendEventWithString(unsigned long long e, const char* info,
                unsigned long long tag);
void UnitySendEventWithNumber(unsigned long long e, unsigned long long info,
                unsigned long long tag);

#endif  // WRAPPER_UNITY_BRIDGE_H

