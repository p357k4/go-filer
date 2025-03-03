Write an application in Go to copy files that:
* propagates errors
* uses logical components and packages
* uses service account key file for gcp authorization
* uses configuration file
* uses logging
* uses interfaces for generalization
* consists of logical components and packages
* each component is run as a separate execution of an application executable with a specific parameter
Monitoring component:
* monitors periodically content of files/landing folder
* moves file without a change in size in two consecutive listings to files/incoming folder
Filtering component:
* analyses content of file in files/incoming folder
* moves files matching filter applied to file content to files/accepted folder
* moves non matching files to files/rejected folder
Uploading component
* upload files from files/accepted folder to gcs bucket
* creates upload confirmation file in files/confirmed folder with: path of source file and with path of target file
* moves successfully uploaded files from files/accepted folder to files/uploaded folder
* use retry after failed upload
* stops trying after several failed upload attempts and move file to files/failed folder
Tracking component:
* monitors files in files/confirmed folder
* reads target file path and adds to manifest file in manifests/incoming folder
* completes manifest file by moving from manifests/incoming to manifests/completed folder after one hour or after reaching 1000 entries
* uploads manifest file from manifests/completed to manifests/incoming gcs folder
* moves uploaded manifest to manifests/uploaded folder
Cleaning component:
* moves uploaded files from manifest file from files/uploaded to files/completed folder
* moves manifest file to manifests/uploaded folder
* use retry after failed upload
* stops trying after several failed upload attempts and move file to manifests/failed folder