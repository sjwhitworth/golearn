#ifndef _H_INTEGRATION_
#define _H_INTEGRATION_

#include "linear.h"

struct problem *CreateCProblem();
void FreeCProblem(struct problem*);
struct model *CreateCModel();
void FreeCModel(struct model*);
struct parameter *CreateCParameter();
void FreeCParameter(struct parameter*);
// Allocates memory outside of golang for describing feature
// vectors. First pointer is the C-managed liblinear problem,
// second parameter is the number of rows we're training on
// third parameter is the total number of non-zero elements 
// (including bias and null terminators) that we need to allocate
// and final parameter is an array describing the number of
// nodes in each row.
int AllocateFeatureNodesForProblem(struct problem*, int, int, int*);
int AllocateLabelsForProblem(struct problem *, int);
void AssignLabelForProblem(struct problem *, int, double);
struct feature_node *GetFeatureNodeForIndex(struct problem *, int, int);

#endif
