// hdt_wrapper.cpp
#include "./libhdt/include/HDTManager.hpp"
#include "./libhdt/include/HDTSpecification.hpp"
#include <string>
#include <iostream>
#include <ctime>

using namespace hdt;

extern "C"
{
    int generateHDTWrapper(const char *filename, const char *baseuri, const char *outputfilename)
    {
        HDTSpecification spec;

        std::string cppFilename(filename);
        // Read RDF into an HDT file.
        HDT *hdt = HDTManager::generateHDT(filename, baseuri, TURTLE, spec);

        // Add additional domain-specific properties to the header
        Header *header = hdt->getHeader();

        time_t rawtime;
        struct tm *timeinfo;
        char buffer[80];

        time(&rawtime);
        timeinfo = localtime(&rawtime);

        strftime(buffer, sizeof(buffer), "%d-%m-%Y %H:%M:%S", timeinfo);
        std::string str(buffer);

        header->insert("http://senforsce.com/o8/Brain/HDTFile", "http://senforsce.com/o8/Brain/LastGenerated", str);

        // Save HDT to a file
        hdt->saveToHDT(outputfilename);

        delete hdt;

        return 0;
    }

    extern "C" int searchWrapper(const char *filename, const char *s, const char *p, const char *o, char *resultBuffer, int bufferSize)
    {
        const int NOT_FOUND = 1;
        const int FOUND = 0;

        int foundResult = NOT_FOUND;
        // Load HDT file
        HDT *hdt = HDTManager::mapHDT(filename);
        // Enumerate all triples matching a pattern ("" means any)
        IteratorTripleString *it = hdt->search(s, p, o);
        unsigned long nb = (unsigned long)(hdt->getTriples()->getNumberOfElements());
        cout << "s:" << s << "p:" << p << "o:" << o << endl;

        string res;
        while (it->hasNext())
        {
            TripleString *triple = it->next();
            cout << "T|" << triple->getSubject() << triple->getPredicate() << triple->getObject() << endl;

            res = triple->getObject();
            foundResult = FOUND;
        }
        delete it;  // Remember to delete iterator to avoid memory leaks!
        delete hdt; // Remember to delete instance when no longer needed!

        strncpy(resultBuffer, res.c_str(), bufferSize);
        resultBuffer[bufferSize - 1] = '\0'; // Null-terminate the string

        return foundResult;
    }
}