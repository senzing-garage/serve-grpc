/**********************************************************************************
 © Copyright Senzing, Inc. 2019-2021
 The source code for this program is not published or otherwise divested
 of its trade secrets, irrespective of what has been deposited with the U.S.
 Copyright Office.
**********************************************************************************/

#ifndef LIBG2DIAGNOSTIC_H
#define LIBG2DIAGNOSTIC_H

/* Platform specific function export header */
#if defined(_WIN32)
  #define _DLEXPORT __declspec(dllexport)
#else
  #include <stddef.h>
  #define _DLEXPORT __attribute__ ((visibility ("default")))
#endif


#ifdef __cplusplus
extern "C"
{
#endif


  /**
   * @brief
   * This method will initialize the G2 diagnostic object.  It must be called
   * prior to any other calls.
   *
   * @param moduleName A name for the diagnostic node, to help identify it within
   *        system logs.
   * @param iniParams A JSON string specifying the configuration parameters
   * @param verboseLogging A flag to enable deeper logging of the G2 processing
   */
  _DLEXPORT int G2Diagnostic_init(const char *moduleName, const char *iniParams, const int verboseLogging);
  _DLEXPORT int G2Diagnostic_initWithConfigID(const char *moduleName, const char *iniParams, const long long initConfigID, const int verboseLogging);


  /**
   * @brief
   * This method will re-initialize the G2 diagnostic object.
   */
  _DLEXPORT int G2Diagnostic_reinit(const long long initConfigID);


  /**
   * @brief
   * This method will destroy and perform cleanup for the G2 diagnostic object.  It
   * should be called after all other calls are complete.
   */
  _DLEXPORT int G2Diagnostic_destroy();


  /**
   * @brief Get the number of physical CPU cores.
   */
  _DLEXPORT int G2Diagnostic_getPhysicalCores();


  /**
   * @brief Get the number of physical CPU cores.
   */
  _DLEXPORT int G2Diagnostic_getLogicalCores();


  /**
   * @brief Get the amount of total system memory
   */
  _DLEXPORT long long G2Diagnostic_getTotalSystemMemory();


  /**
   * @brief Get the amount of available system memory
   */
  _DLEXPORT long long G2Diagnostic_getAvailableMemory();


  /**
   * @brief Check the performance metrics of the datastore
   */
  _DLEXPORT int G2Diagnostic_checkDBPerf(int secondsToRun, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize) );


  /**
   * @brief Get information for the datastore
   */
  _DLEXPORT int G2Diagnostic_getDBInfo(char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize) );


  /**
   * @brief Retrieve diagnostic information on the contents of the data store
   */
  _DLEXPORT int G2Diagnostic_getDataSourceCounts(char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize) );
  _DLEXPORT int G2Diagnostic_getMappingStatistics(const int includeInternalFeatures, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize) );
  _DLEXPORT int G2Diagnostic_getGenericFeatures(const char* featureType, const size_t maximumEstimatedCount, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize) );
  _DLEXPORT int G2Diagnostic_getEntitySizeBreakdown(const size_t minimumEntitySize, const int includeInternalFeatures, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize) );
  _DLEXPORT int G2Diagnostic_getEntityDetails(const long long entityID, const int includeInternalFeatures, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize) );
  _DLEXPORT int G2Diagnostic_getResolutionStatistics(char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize) );
  _DLEXPORT int G2Diagnostic_getRelationshipDetails(const long long relationshipID, const int includeInternalFeatures, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize) );
  _DLEXPORT int G2Diagnostic_getEntityResume(const long long entityID, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize) );


  /**
   * @brief Retrieve diagnostic information on sized-entities in the data store
   */
  typedef void* EntityListBySizeHandle;
  _DLEXPORT int G2Diagnostic_getEntityListBySize(const size_t entitySize,EntityListBySizeHandle* entityListBySizeHandle);
  _DLEXPORT int G2Diagnostic_fetchNextEntityBySize(EntityListBySizeHandle entityListBySizeHandle, char *responseBuf, const size_t bufSize);
  _DLEXPORT int G2Diagnostic_closeEntityListBySize(EntityListBySizeHandle entityListBySizeHandle);


  /**
   * @brief
   * This function asks g2 for any entities having any of the lib feat id specified in the "features" doc.
   * The "features" also contains an entity id. This entity is ignored in the returned values.
   *
   * @param features Json document of the format: {"ENTITY_ID":<entity id>,"LIB_FEAT_IDS":[<id1>,<id2>,...<idn>]}
   * where entity id specifies the entity to ignore in the returns and <id#> are the lib feat ids used to query for
   * entities.
   * @param responseBuf Json document in the format:
   * [{"LIB_FEAT_ID":<lib feat id>, "USAGE_TYPE":"<usage type","RES_ENT_ID":<entity id1>},...]
   * @param bufSize Size of responseBuff in bytes
   * @param resizeFunc Function for resizing the responseBuf it it exeeds the original size.
   *
   * @return 0 on successful and response is a JSON string of the detail
   */
  _DLEXPORT int G2Diagnostic_findEntitiesByFeatureIDs(const char *features, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize));

  /**
   * @brief
   * This method is used to retrieve a stored feature.
   * 
   * @param libFeatID The ID of the feature to search for
   * @param responseBuf A memory buffer for returning the response document.
   *        If an error occurred, an error response is stored here.
   * @param bufSize The max number of bytes that can be stored in responseBuf.
   *        The response buffer must be able to hold at least this many bytes,
   *        or a resize method must be provided that may be used to resize
   *        the buffer.  This will return the new size.
   * @param resizeFunc A function pointer that can be used to resize the memory
   *        buffer specified in the responseBuf argument.  This function will
   *        be called to allocate more memory if the response buffer is not large
   *        enough.  This argument may be NULL.  If it is NULL, the function
   *        will return an error if the result is larger than the buffer.
   * @return Returns 0 for success. Returns -1 if the response status indicates
   *         failure or the G2 module is not initialized. Returns -2 if 
   *         an exception was thrown and caught.
   */
  _DLEXPORT int G2Diagnostic_getFeature(const long long libFeatID, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize));


  /**
   * @brief
   * This function retrieves the last exception thrown in G2Diagnostic
   * @return number of bytes copied into buffer
   */
  _DLEXPORT int G2Diagnostic_getLastException(char *buffer, const size_t bufSize);


  /**
   * @brief
   * This function retrieves the code of the last exception thrown in G2Diagnostic
   * @return number of bytes copied into buffer
   */
  _DLEXPORT int G2Diagnostic_getLastExceptionCode();


  /**
   * @brief
   * This function clears the last exception thrown in G2Diagnostic
   */
  _DLEXPORT void G2Diagnostic_clearLastException();


#ifdef __cplusplus
};
#endif

#endif /* LIBG2DIAGNOSTIC_H */
